package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/common/utils"

	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCountByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCountByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCountByIdLogic {
	return &GetCountByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCountByIdLogic) GetCountById(in *video.GetCountByIdRequest) (*video.GetCountByIdResponse, error) {
	db := l.svcCtx.DB
	var workCount, totalFavorited, userFavoriteCount int64
	err := db.Model(&model.Video{}).
		Where("author_id = ?", in.UserId).
		Count(&workCount).Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying WorkCount ", err)
	}
	err = db.Model(&model.Favorite{}).
		Where("user_id = ?", in.UserId).
		Count(&userFavoriteCount).Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying UserFavoriteCount ", err)
	}

	err = db.Model(&model.Favorite{}).
		Where("video_id in (?)",
			db.Model(&model.Video{}).
				Select("video_id").
				Where("author_id = ?", in.UserId)).
		Count(&totalFavorited).Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying UserFavoriteCount ", err)
	}
	return &video.GetCountByIdResponse{
		TotalFavorited:    totalFavorited,
		WorkCount:         workCount,
		UserFavoriteCount: userFavoriteCount,
	}, nil
}
