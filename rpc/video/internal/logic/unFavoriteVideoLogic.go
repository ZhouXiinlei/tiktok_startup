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

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&model.Video{}).
			Where("video_id = ?", in.VideoId).
			Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).
			Error
		if err != nil {
			return utils.InternalWithDetails("error reducing favorite_count", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &video.Empty{}, nil
}
