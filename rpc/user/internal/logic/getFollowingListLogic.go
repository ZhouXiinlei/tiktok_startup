package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/common/utils"

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
	userList := make([]*model.User, 0)
	err := l.svcCtx.DB.
		Where("user_id in (?)", l.svcCtx.DB.
			Model(&model.Follow{}).
			Select("followed_id").
			Where("follower_id = ?", in.UserId)).
		Find(&userList).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying following list", err)
	}

	followingList := make([]*user.UserInfo, 0, len(userList))
	for _, followed := range userList {
		followingList = append(followingList, &user.UserInfo{
			UserId:         followed.UserId,
			Username:       followed.Username,
			FollowingCount: followed.FollowingCount,
			FollowerCount:  followed.FollowerCount,
			CreatedAt:      followed.CreatedAt.Unix(),
			UpdatedAt:      followed.UpdatedAt.Unix(),
			IsFollow:       true,
		})
	}
	return &user.GetFollowingListResponse{
		FollowingList: followingList,
	}, nil
}
