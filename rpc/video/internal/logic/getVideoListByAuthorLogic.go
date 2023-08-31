package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/rpc/video/internal/svc"
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
	err := l.svcCtx.Mysql.Where("author_id = ?", in.AuthorId).
		Order("created_at desc").
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	resp := &video.GetVideoListByAuthorResponse{}
	for _, v := range videos {
		videoInfo := &video.VideoInfo{
			Id:            int64(v.ID),
			AuthorId:      v.AuthorId,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
		}
		resp.Video = append(resp.Video, videoInfo)
	}
	return resp, nil
}
