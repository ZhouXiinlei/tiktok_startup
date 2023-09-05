package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/internal/cache"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"
)

type UnFavoriteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFavoriteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFavoriteVideoLogic {
	return &UnFavoriteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFavoriteVideoLogic) UnFavoriteVideo(in *video.UnFavoriteVideoRequest) (*video.Empty, error) {
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND video_id = ?", in.UserId, in.VideoId).Delete(&model.Favorite{})
		if err := res.Error; err != nil {
			return utils.InternalWithDetails("error deleting favorite record", err)
		}
		if res.RowsAffected == 0 {
			return nil
		}
		err := l.svcCtx.RDS.Set(cache.GenFavoriteKey(in.UserId, in.VideoId), "no")
		if err != nil {
			return utils.InternalWithDetails("(redis)error updating favorite relation", err)
		}
		err = cache.ModifyVideoCounts(l.svcCtx.DB, l.svcCtx.RDS, in.VideoId, "favorite_count", -1)
		if err != nil {
			return utils.InternalWithDetails("error adding favorite_count", err)
		}
		var targetId int64 = 0
		err = l.svcCtx.DB.Model(&model.Video{}).Where("video_id = ?", in.VideoId).Select("author_id").First(&targetId).Error
		if err != nil {
			return utils.InternalWithDetails("error querying video_author record", err)
		}
		_, err = l.svcCtx.UserRpc.ModFavorite(l.ctx, &userClient.ModFavoriteRequest{
			UserId:   in.UserId,
			TargetId: targetId,
			Delta:    -1,
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
