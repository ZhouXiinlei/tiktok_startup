package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/rpc/user/internal/union"

	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowingListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowingListLogic {
	return &GetFollowingListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowingListLogic) GetFollowingList(in *user.GetFollowingListRequest) (*user.GetFollowingListResponse, error) {
	userInfoList, err := union.GetManyUserInfos(l.svcCtx.DB, l.svcCtx.RDS, l.svcCtx.DB.
		Model(&model.Follow{}).
		Select("followed_id").
		Where("follower_id = ?", in.UserId))
	if err != nil {
		return nil, err
	}

	for _, item := range userInfoList {
		item.IsFollow = true
	}

	return &user.GetFollowingListResponse{
		FollowingList: userInfoList,
	}, nil
}
