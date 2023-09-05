package social

import (
	"context"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowListLogic) GetFollowList(req *types.GetFollowListRequest) (resp *types.GetFollowListResponse, err error) {
	GetFollowListData, err := l.svcCtx.UserRpc.GetFollowingList(l.ctx, &user.GetFollowingListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetFollowList failed, err:%v", err)
		return
	}
	var followList []types.User
	for _, follow := range GetFollowListData.FollowingList {
		followList = append(followList, types.User{
			Id:             follow.UserId,
			Name:           follow.Username,
			FollowCount:    follow.FollowingCount,
			FollowerCount:  follow.FollowerCount,
			IsFollow:       follow.IsFollow,
			TotalFavorited: follow.TotalFavorited,
			WorkCount:      follow.WorkCount,
			FavoriteCount:  follow.FavoriteCount,
		})
	}
	return &types.GetFollowListResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: followList,
	}, nil
}
