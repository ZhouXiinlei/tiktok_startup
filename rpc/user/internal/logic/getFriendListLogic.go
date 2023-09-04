package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/common/utils"

	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *user.GetFriendListRequest) (*user.GetFriendListResponse, error) {
	var count int64
	err := l.svcCtx.DB.Model(&model.User{}).Where("user_id = ?", in.UserId).Count(&count).Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying user record", err)
	}
	if count == 0 {
		return nil, common.ErrUserNotFound.Err()
	}

	recordList := make([]*model.Friend, 0)
	err = l.svcCtx.DB.
		Where("user_a_id = ?", in.UserId).
		Or("user_b_id = ?", in.UserId).
		Order("created_at DESC").
		Find(&recordList).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying friend record", err)
	}

	order := make(map[int64]int, len(recordList))
	friendList, err := mr.MapReduce(func(source chan<- *model.Friend) {
		for i, friendRecord := range recordList {
			source <- friendRecord
			order[friendRecord.FriendId] = i
		}
	}, func(friendRecord *model.Friend, writer mr.Writer[*FriendListPipePack], cancel func(error)) {
		var friendId int64
		if friendRecord.UserAId == in.UserId {
			friendId = friendRecord.UserBId
		} else {
			friendId = friendRecord.UserAId
		}

		userRecord := &model.User{}
		err := l.svcCtx.DB.Where("user_id = ?", friendId).First(userRecord).Error
		if err != nil {
			cancel(utils.InternalWithDetails("error querying User record", err))
			return
		}

		writer.Write(&FriendListPipePack{
			User:   userRecord,
			Friend: friendRecord,
		})
	}, func(pipe <-chan *FriendListPipePack, writer mr.Writer[[]*user.UserInfo], cancel func(error)) {
		list := make([]*user.UserInfo, len(recordList))

		for pack := range pipe {
			i, _ := order[pack.Friend.FriendId]
			list[i] = &user.UserInfo{
				UserId:         pack.User.UserId,
				Username:       pack.User.Username,
				FollowingCount: pack.User.FollowingCount,
				FollowerCount:  pack.User.FollowerCount,
				IsFollow:       true,
				CreatedAt:      pack.User.CreatedAt.Unix(),
				UpdatedAt:      pack.User.UpdatedAt.Unix(),
			}
		}

		writer.Write(list)
	})
	if err != nil {
		return nil, err
	}

	return &user.GetFriendListResponse{
		FriendList: friendList,
	}, nil
}

type FriendListPipePack struct {
	User   *model.User
	Friend *model.Friend
}
