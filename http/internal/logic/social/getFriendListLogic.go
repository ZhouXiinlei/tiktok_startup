package social

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"tikstart/common"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/user/user"
)

type GetFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListRequest) (resp *types.GetFriendListResponse, err error) {
	queryId := req.UserId
	friendListRes, err := l.svcCtx.UserRpc.GetFriendList(l.ctx, &user.GetFriendListRequest{
		UserId: queryId,
	})
	if err != nil {
		if st, match := utils.MatchError(err, common.ErrUserNotFound); match {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42202,
				Message:    "用户不存在",
			}
		} else {
			logx.WithContext(l.ctx).Errorf("获取用户好友列表失败: %v", err)
			return nil, utils.ReturnInternalError(st, err)
		}
	}

	userList := make([]types.User, 0, len(friendListRes.FriendList))
	for _, friend := range friendListRes.FriendList {
		//countRes, err := l.svcCtx.VideoRpc.GetCountById(l.ctx, &video.GetCountByIdRequest{
		//	UserId: friend.UserId,
		//})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户视频统计数据失败: %v", err)
			return nil, utils.ReturnInternalError(status.Convert(err), err)
		}

		userList = append(userList, types.User{
			Id:             friend.UserId,
			Name:           friend.Username,
			FollowCount:    friend.FollowingCount,
			FollowerCount:  friend.FollowerCount,
			IsFollow:       friend.IsFollow,
			TotalFavorited: friend.TotalFavorited,
			WorkCount:      friend.WorkCount,
			FavoriteCount:  friend.FavoriteCount,
		})
	}

	return &types.GetFriendListResponse{
		BasicResponse: types.BasicResponse{},
		UserList:      userList,
	}, nil
}
