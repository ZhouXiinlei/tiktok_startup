// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package userClient

import (
	"context"

	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateRequest            = user.CreateRequest
	CreateResponse           = user.CreateResponse
	Empty                    = user.Empty
	FollowRequest            = user.FollowRequest
	GetFollowerListRequest   = user.GetFollowerListRequest
	GetFollowerListResponse  = user.GetFollowerListResponse
	GetFollowingListRequest  = user.GetFollowingListRequest
	GetFollowingListResponse = user.GetFollowingListResponse
	GetFriendListRequest     = user.GetFriendListRequest
	GetFriendListResponse    = user.GetFriendListResponse
	IsFollowRequest          = user.IsFollowRequest
	IsFollowResponse         = user.IsFollowResponse
	PingRequest              = user.PingRequest
	PingResponse             = user.PingResponse
	QueryByIdRequest         = user.QueryByIdRequest
	QueryByNameRequest       = user.QueryByNameRequest
	QueryResponse            = user.QueryResponse
	UnFollowRequest          = user.UnFollowRequest
	UserInfo                 = user.UserInfo

	User interface {
		Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
		Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
		QueryById(ctx context.Context, in *QueryByIdRequest, opts ...grpc.CallOption) (*QueryResponse, error)
		QueryByName(ctx context.Context, in *QueryByNameRequest, opts ...grpc.CallOption) (*QueryResponse, error)
		Follow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*Empty, error)
		UnFollow(ctx context.Context, in *UnFollowRequest, opts ...grpc.CallOption) (*Empty, error)
		GetFollowerList(ctx context.Context, in *GetFollowerListRequest, opts ...grpc.CallOption) (*GetFollowerListResponse, error)
		GetFollowingList(ctx context.Context, in *GetFollowingListRequest, opts ...grpc.CallOption) (*GetFollowingListResponse, error)
		IsFollow(ctx context.Context, in *IsFollowRequest, opts ...grpc.CallOption) (*IsFollowResponse, error)
		GetFriendList(ctx context.Context, in *GetFriendListRequest, opts ...grpc.CallOption) (*GetFriendListResponse, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

func (m *defaultUser) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}

func (m *defaultUser) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Create(ctx, in, opts...)
}

func (m *defaultUser) QueryById(ctx context.Context, in *QueryByIdRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.QueryById(ctx, in, opts...)
}

func (m *defaultUser) QueryByName(ctx context.Context, in *QueryByNameRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.QueryByName(ctx, in, opts...)
}

func (m *defaultUser) Follow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*Empty, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Follow(ctx, in, opts...)
}

func (m *defaultUser) UnFollow(ctx context.Context, in *UnFollowRequest, opts ...grpc.CallOption) (*Empty, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.UnFollow(ctx, in, opts...)
}

func (m *defaultUser) GetFollowerList(ctx context.Context, in *GetFollowerListRequest, opts ...grpc.CallOption) (*GetFollowerListResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.GetFollowerList(ctx, in, opts...)
}

func (m *defaultUser) GetFollowingList(ctx context.Context, in *GetFollowingListRequest, opts ...grpc.CallOption) (*GetFollowingListResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.GetFollowingList(ctx, in, opts...)
}

func (m *defaultUser) IsFollow(ctx context.Context, in *IsFollowRequest, opts ...grpc.CallOption) (*IsFollowResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.IsFollow(ctx, in, opts...)
}

func (m *defaultUser) GetFriendList(ctx context.Context, in *GetFriendListRequest, opts ...grpc.CallOption) (*GetFriendListResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.GetFriendList(ctx, in, opts...)
}
