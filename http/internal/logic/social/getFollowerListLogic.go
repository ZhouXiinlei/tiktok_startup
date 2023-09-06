package social

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"tikstart/rpc/user/user"

	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerListLogic {
	return &GetFollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowerListLogic) GetFollowerList(req *types.GetFollowerListRequest) (resp *types.GetFollowerListResponse, err error) {
	GetFollowerListData, err := l.svcCtx.UserRpc.GetFollowerList(l.ctx, &user.GetFollowerListRequest{
		UserId: req.UserId,
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetFollowerList failed, err:%v", err)
		return
	}

	order := make(map[int]int, len(GetFollowerListData.FollowerList))

	followerList, err := mr.MapReduce(func(source chan<- interface{}) {
		for i, v := range GetFollowerListData.FollowerList {
			source <- v
			order[int(v.UserId)] = i
		}
	}, func(item interface{}, writer mr.Writer[types.User], cancel func(error)) {
		follower := item.(*user.UserInfo)

		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:   req.UserId,
			TargetId: follower.UserId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("IsFollow failed, err:%v", err)
			cancel(err)
			return
		}
		writer.Write(types.User{
			Id:             follower.UserId,
			Name:           follower.Username,
			FollowCount:    follower.FollowingCount,
			FollowerCount:  follower.FollowerCount,
			IsFollow:       isFollowRes.IsFollow,
			TotalFavorited: follower.TotalFavorited,
			WorkCount:      follower.WorkCount,
			FavoriteCount:  follower.FavoriteCount,
		})
	}, func(pipe <-chan types.User, writer mr.Writer[[]types.User], cancel func(error)) {
		list := make([]types.User, len(GetFollowerListData.FollowerList))
		for item := range pipe {
			temp := item
			i, ok := order[int(temp.Id)]
			if !ok {
				cancel(err)
				return
			}
			list[i] = temp
		}
		writer.Write(list)
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetFollowerList failed, err:%v", err)
		return
	}

	return &types.GetFollowerListResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: followerList,
	}, nil
}
