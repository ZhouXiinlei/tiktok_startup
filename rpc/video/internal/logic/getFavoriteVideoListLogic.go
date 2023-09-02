package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"
)

type GetFavoriteVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteVideoListLogic {
	return &GetFavoriteVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteVideoListLogic) GetFavoriteVideoList(in *video.GetFavoriteVideoListRequest) (*video.GetFavoriteVideoListResponse, error) {
	var favoriteVideoList []*model.Favorite

	if err := l.svcCtx.DB.
		Where("user_id = ?", in.UserId).
		Preload("Video").
		Order("created_at desc").
		Find(&favoriteVideoList).
		Error; err != nil {
		return nil, utils.InternalWithDetails("error querying favorite video list", err)
	}

	videoList := make([]*video.VideoInfo, 0, len(favoriteVideoList))
	for _, v := range favoriteVideoList {
		videoInfo := &video.VideoInfo{
			Id:            v.Video.VideoId,
			AuthorId:      v.Video.AuthorId,
			Title:         v.Video.Title,
			PlayUrl:       v.Video.PlayUrl,
			CoverUrl:      v.Video.CoverUrl,
			FavoriteCount: v.Video.FavoriteCount,
			CommentCount:  v.Video.CommentCount,
		}
		videoList = append(videoList, videoInfo)
	}
	return &video.GetFavoriteVideoListResponse{
		VideoList: videoList,
	}, nil
}
