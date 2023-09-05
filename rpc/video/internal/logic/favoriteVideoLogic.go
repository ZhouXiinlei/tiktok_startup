package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/internal/cache"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"
)

type FavoriteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteVideoLogic {
	return &FavoriteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteVideoLogic) FavoriteVideo(in *video.FavoriteVideoRequest) (*video.Empty, error) {
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		err := tx.Model(&model.Video{}).Where("video_id = ?", in.VideoId).Count(&count).Error
		if err != nil {
			return utils.InternalWithDetails("error querying video record", err)
		}
		if count == 0 {
			return common.ErrVideoNotFound.Err()
		}

		res, err := cache.IsFavorite(l.svcCtx, in.UserId, in.VideoId)
		if err != nil {
			return err
		}
		if res {
			return nil
		}

		err = tx.Create(&model.Favorite{
			UserId:  in.UserId,
			VideoId: in.VideoId,
		}).Error
		if err != nil {
			return utils.InternalWithDetails("error creating favorite record", err)
		}

		err = l.svcCtx.RDS.Set(cache.GenFavoriteKey(in.UserId, in.VideoId), "yes")
		if err != nil {
			return utils.InternalWithDetails("(redis)error updating favorite relation", err)
		}
		err = cache.ModifyVideoCounts(tx, l.svcCtx.RDS, in.VideoId, "favorite_count", 1)
		if err != nil {
			return utils.InternalWithDetails("error adding favorite_count", err)
		}
		var targetId int64 = 0
		err = tx.Model(&model.Video{}).Where("video_id = ?", in.VideoId).Select("author_id").First(&targetId).Error
		if err != nil {
			return utils.InternalWithDetails("error querying video_author record", err)
		}
		_, err = l.svcCtx.UserRpc.ModFavorite(l.ctx, &userClient.ModFavoriteRequest{
			UserId:   in.UserId,
			TargetId: targetId,
			Delta:    1,
		})
		if err != nil {
			return utils.InternalWithDetails("error updating user_favorite_count", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &video.Empty{}, nil
}
