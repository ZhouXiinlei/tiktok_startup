package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"tikstart/common/model"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"
)

type PublishVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVideoLogic {
	return &PublishVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishVideoLogic) PublishVideo(in *video.PublishVideoRequest) (*video.Empty, error) {
	newVideo := &model.Video{
		AuthorId: in.Video.AuthorId,
		Title:    in.Video.Title,
		PlayUrl:  in.Video.PlayUrl,
		CoverUrl: in.Video.CoverUrl,
	}
	if err := l.svcCtx.DB.Create(newVideo).Error; err != nil {
		return nil, status.Error(1000, err.Error())
	}
	_, err := l.svcCtx.UserRpc.ModWorkCount(l.ctx, &userClient.ModWorkCountRequest{
		UserId: in.Video.AuthorId,
		Delta:  1,
	})
	if err != nil {
		return nil, err
	}
	return &video.Empty{}, nil
}
