package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/internal/union"
	"tikstart/rpc/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListByAuthorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListByAuthorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListByAuthorLogic {
	return &GetVideoListByAuthorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoListByAuthorLogic) GetVideoListByAuthor(in *video.GetVideoListByAuthorRequest) (*video.GetVideoListByAuthorResponse, error) {
	var videos []model.Video
	err := l.svcCtx.DB.
		Where("author_id = ?", in.AuthorId).
		Order("created_at desc").
		Find(&videos).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("error getting video list by author", err)
	}

	videoList := make([]*video.VideoInfo, 0, len(videos))
	for _, v := range videos {
		videoInfo, err := union.GetVideoInfoById(l.svcCtx.DB, l.svcCtx.RDS, v.VideoId)
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, videoInfo)
	}
	return &video.GetVideoListByAuthorResponse{
		Video: videoList,
	}, nil
}
