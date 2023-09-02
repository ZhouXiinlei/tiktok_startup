package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"
)

type IsFavoriteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFavoriteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFavoriteVideoLogic {
	return &IsFavoriteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFavoriteVideoLogic) IsFavoriteVideo(in *video.IsFavoriteVideoRequest) (*video.IsFavoriteVideoResponse, error) {
	var count int64
	err := l.svcCtx.Mysql.
		Model(&model.Favorite{}).
		Where("user_id = ? AND video_id = ?", in.UserId, in.VideoId).
		Count(&count).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying favorite count", err)
	}

	if count == 0 {
		return &video.IsFavoriteVideoResponse{
			IsFavorite: false,
		}, nil
	}
	return &video.IsFavoriteVideoResponse{
		IsFavorite: true,
	}, nil
}
