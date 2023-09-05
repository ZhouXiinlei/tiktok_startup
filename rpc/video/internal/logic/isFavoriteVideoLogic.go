package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/rpc/video/internal/cache"
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
	res, err := cache.IsFavorite(l.svcCtx, in.UserId, in.VideoId)
	if err != nil {
		logx.WithContext(l.ctx).Error(err)
		return nil, err
	}

	return &video.IsFavoriteVideoResponse{
		IsFavorite: res,
	}, nil
}
