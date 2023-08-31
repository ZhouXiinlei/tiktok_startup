package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"
)

type GetFollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerListLogic {
	return &GetFollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerListLogic) GetFollowerList(in *user.GetFollowerListRequest) (*user.GetFollowerListResponse, error) {
	userPreload := &model.User{}
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).Preload("Followers").Find(&userPreload).Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying follower list", err)
	}

	followerList := make([]*user.UserInfo, 0, len(userPreload.Followers))
	for _, follower := range userPreload.Followers {
		var count int64
		err := l.svcCtx.DB.Model(&model.Follow{}).Where("follower_id = ? AND followed_id = ?", in.UserId, follower.UserId).Count(&count).Error
		if err != nil {
			return nil, utils.InternalWithDetails("error querying follow relation", err)
		}

		followerList = append(followerList, &user.UserInfo{
			UserId:         follower.UserId,
			Username:       follower.Username,
			FollowingCount: follower.FollowingCount,
			FollowerCount:  follower.FollowerCount,
			CreatedAt:      follower.CreatedAt.Unix(),
			UpdatedAt:      follower.UpdatedAt.Unix(),
			IsFollow:       count == 1,
		})
	}
	return &user.GetFollowerListResponse{
		FollowerList: followerList,
	}, nil
}
