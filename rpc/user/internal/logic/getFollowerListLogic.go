package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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
	userList := make([]*model.User, 0)
	err := l.svcCtx.DB.
		Where("user_id in (?)", l.svcCtx.DB.
			Model(&model.Follow{}).
			Select("follower_id").
			Where("followed_id = ?", in.UserId)).
		Find(&userList).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying follower list", err)
	}

	order := make(map[int]int, len(userList))

	followerList, err := mr.MapReduce(func(source chan<- *model.User) {
		for i, follower := range userList {
			source <- follower
			order[int(follower.UserId)] = i
		}
	}, func(item *model.User, writer mr.Writer[*user.UserInfo], cancel func(error)) {
		var count int64
		err := l.svcCtx.DB.Model(&model.Follow{}).Where("follower_id = ? AND followed_id = ?", in.UserId, item.UserId).Count(&count).Error
		if err != nil {
			cancel(utils.InternalWithDetails("error querying follow relation", err))
			return
		}

		writer.Write(&user.UserInfo{
			UserId:         item.UserId,
			Username:       item.Username,
			FollowingCount: item.FollowingCount,
			FollowerCount:  item.FollowerCount,
			CreatedAt:      item.CreatedAt.Unix(),
			UpdatedAt:      item.UpdatedAt.Unix(),
			IsFollow:       count == 1,
		})

	}, func(pipe <-chan *user.UserInfo, writer mr.Writer[[]*user.UserInfo], cancel func(error)) {
		list := make([]*user.UserInfo, len(userList))
		for item := range pipe {
			userInfo := item
			i, _ := order[int(item.UserId)]
			list[i] = userInfo
		}
		writer.Write(list)
	})

	return &user.GetFollowerListResponse{
		FollowerList: followerList,
	}, nil
}
