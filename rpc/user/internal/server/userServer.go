// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"tikstart/rpc/user/internal/logic"
	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Ping(ctx context.Context, in *user.PingRequest) (*user.PingResponse, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}

func (s *UserServer) Create(ctx context.Context, in *user.CreateRequest) (*user.CreateResponse, error) {
	l := logic.NewCreateLogic(ctx, s.svcCtx)
	return l.Create(in)
}

func (s *UserServer) QueryById(ctx context.Context, in *user.QueryByIdRequest) (*user.UserInfo, error) {
	l := logic.NewQueryByIdLogic(ctx, s.svcCtx)
	return l.QueryById(in)
}

func (s *UserServer) QueryByName(ctx context.Context, in *user.QueryByNameRequest) (*user.UserInfo, error) {
	l := logic.NewQueryByNameLogic(ctx, s.svcCtx)
	return l.QueryByName(in)
}

func (s *UserServer) Follow(ctx context.Context, in *user.FollowRequest) (*user.Empty, error) {
	l := logic.NewFollowLogic(ctx, s.svcCtx)
	return l.Follow(in)
}

func (s *UserServer) UnFollow(ctx context.Context, in *user.UnFollowRequest) (*user.Empty, error) {
	l := logic.NewUnFollowLogic(ctx, s.svcCtx)
	return l.UnFollow(in)
}

func (s *UserServer) GetFollowerList(ctx context.Context, in *user.GetFollowerListRequest) (*user.GetFollowerListResponse, error) {
	l := logic.NewGetFollowerListLogic(ctx, s.svcCtx)
	return l.GetFollowerList(in)
}

func (s *UserServer) GetFollowingList(ctx context.Context, in *user.GetFollowingListRequest) (*user.GetFollowingListResponse, error) {
	l := logic.NewGetFollowingListLogic(ctx, s.svcCtx)
	return l.GetFollowingList(in)
}

func (s *UserServer) IsFollow(ctx context.Context, in *user.IsFollowRequest) (*user.IsFollowResponse, error) {
	l := logic.NewIsFollowLogic(ctx, s.svcCtx)
	return l.IsFollow(in)
}

func (s *UserServer) GetFriendList(ctx context.Context, in *user.GetFriendListRequest) (*user.GetFriendListResponse, error) {
	l := logic.NewGetFriendListLogic(ctx, s.svcCtx)
	return l.GetFriendList(in)
}

func (s *UserServer) ModFavorite(ctx context.Context, in *user.ModFavoriteRequest) (*user.Empty, error) {
	l := logic.NewModFavoriteLogic(ctx, s.svcCtx)
	return l.ModFavorite(in)
}

func (s *UserServer) ModWorkCount(ctx context.Context, in *user.ModWorkCountRequest) (*user.Empty, error) {
	l := logic.NewModWorkCountLogic(ctx, s.svcCtx)
	return l.ModWorkCount(in)
}
