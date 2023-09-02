package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tikstart/common/model"
	"tikstart/common/utils"
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
	err := l.svcCtx.Mysql.Transaction(func(tx *gorm.DB) error {
		var count int64
		err := tx.
			Model(&model.Favorite{}).
			Where("user_id = ? AND video_id = ?", in.UserId, in.VideoId).
			Count(&count).
			Error
		if err != nil {
			return utils.InternalWithDetails("error querying favorite record", err)
		}
		if count > 0 {
			return nil
		}

		err = tx.Create(&model.Favorite{
			UserId:  in.UserId,
			VideoId: in.VideoId,
		}).Error
		if err != nil {
			return utils.InternalWithDetails("error creating favorite record", err)
		}

		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&model.Video{}).
			Where("video_id = ?", in.VideoId).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).
			Error
		if err != nil {
			return utils.InternalWithDetails("error adding favorite count", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &video.Empty{}, nil
}
