package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"tikstart/common/model"
	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/internal/union"
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
	userInfoList, err := union.GetManyUserInfos(l.svcCtx.DB, l.svcCtx.RDS, l.svcCtx.DB.
		Model(&model.Follow{}).
		Select("follower_id").
		Where("followed_id = ?", in.UserId))
	if err != nil {
		return nil, err
	}

	order := make(map[int64]int, len(userInfoList))

	followerList, err := mr.MapReduce(func(source chan<- *user.UserInfo) {
		for i, follower := range userInfoList {
			source <- follower
			order[follower.UserId] = i
		}
	}, func(item *user.UserInfo, writer mr.Writer[*user.UserInfo], cancel func(error)) {
		res, err := union.IsFollow(l.svcCtx.DB, l.svcCtx.RDS, in.UserId, item.UserId)
		if err != nil {
			cancel(err)
			return
		}
		item.IsFollow = res

		writer.Write(item)
	}, func(pipe <-chan *user.UserInfo, writer mr.Writer[[]*user.UserInfo], cancel func(error)) {
		list := make([]*user.UserInfo, len(userInfoList))
		for item := range pipe {
			i, _ := order[item.UserId]
			list[i] = item
		}
		writer.Write(list)
	})

	return &user.GetFollowerListResponse{
		FollowerList: followerList,
	}, nil
}
