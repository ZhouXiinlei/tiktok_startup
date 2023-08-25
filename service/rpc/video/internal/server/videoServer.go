// Code generated by goctl. DO NOT EDIT.
// Source: video.proto

package server

import (
	"context"

	"tiktok_startup/service/rpc/video/internal/logic"
	"tiktok_startup/service/rpc/video/internal/svc"
	"tiktok_startup/service/rpc/video/video"
)

type VideoServer struct {
	svcCtx *svc.ServiceContext
	video.UnimplementedVideoServer
}

func NewVideoServer(svcCtx *svc.ServiceContext) *VideoServer {
	return &VideoServer{
		svcCtx: svcCtx,
	}
}

func (s *VideoServer) GetVideoList(ctx context.Context, in *video.GetVideoListRequest) (*video.GetVideoListResponse, error) {
	l := logic.NewGetVideoListLogic(ctx, s.svcCtx)
	return l.GetVideoList(in)
}

func (s *VideoServer) PublishVideo(ctx context.Context, in *video.PublishVideoRequest) (*video.Empty, error) {
	l := logic.NewPublishVideoLogic(ctx, s.svcCtx)
	return l.PublishVideo(in)
}

func (s *VideoServer) UpdateVideo(ctx context.Context, in *video.UpdateVideoRequest) (*video.Empty, error) {
	l := logic.NewUpdateVideoLogic(ctx, s.svcCtx)
	return l.UpdateVideo(in)
}

func (s *VideoServer) GetVideoListByAuthor(ctx context.Context, in *video.GetVideoListByAuthorRequest) (*video.GetVideoListByAuthorResponse, error) {
	l := logic.NewGetVideoListByAuthorLogic(ctx, s.svcCtx)
	return l.GetVideoListByAuthor(in)
}

func (s *VideoServer) FavoriteVideo(ctx context.Context, in *video.FavoriteVideoRequest) (*video.Empty, error) {
	l := logic.NewFavoriteVideoLogic(ctx, s.svcCtx)
	return l.FavoriteVideo(in)
}

func (s *VideoServer) UnFavoriteVideo(ctx context.Context, in *video.UnFavoriteVideoRequest) (*video.Empty, error) {
	l := logic.NewUnFavoriteVideoLogic(ctx, s.svcCtx)
	return l.UnFavoriteVideo(in)
}

func (s *VideoServer) GetFavoriteVideoList(ctx context.Context, in *video.GetFavoriteVideoListRequest) (*video.GetFavoriteVideoListResponse, error) {
	l := logic.NewGetFavoriteVideoListLogic(ctx, s.svcCtx)
	return l.GetFavoriteVideoList(in)
}

func (s *VideoServer) IsFavoriteVideo(ctx context.Context, in *video.IsFavoriteVideoRequest) (*video.IsFavoriteVideoResponse, error) {
	l := logic.NewIsFavoriteVideoLogic(ctx, s.svcCtx)
	return l.IsFavoriteVideo(in)
}

func (s *VideoServer) CommentVideo(ctx context.Context, in *video.CommentVideoRequest) (*video.CommentVideoResponse, error) {
	l := logic.NewCommentVideoLogic(ctx, s.svcCtx)
	return l.CommentVideo(in)
}

func (s *VideoServer) GetCommentList(ctx context.Context, in *video.GetCommentListRequest) (*video.GetCommentListResponse, error) {
	l := logic.NewGetCommentListLogic(ctx, s.svcCtx)
	return l.GetCommentList(in)
}

func (s *VideoServer) DeleteVideoComment(ctx context.Context, in *video.DeleteVideoCommentRequest) (*video.Empty, error) {
	l := logic.NewDeleteVideoCommentLogic(ctx, s.svcCtx)
	return l.DeleteVideoComment(in)
}

func (s *VideoServer) GetCommentInfo(ctx context.Context, in *video.GetCommentInfoRequest) (*video.GetCommentInfoResponse, error) {
	l := logic.NewGetCommentInfoLogic(ctx, s.svcCtx)
	return l.GetCommentInfo(in)
}
