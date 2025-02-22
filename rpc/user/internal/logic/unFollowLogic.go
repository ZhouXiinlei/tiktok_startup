package logic

import (
	"context"
	"gorm.io/gorm"
	"tikstart/common/cache"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/union"

	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFollowLogic) UnFollow(in *user.UnFollowRequest) (*user.Empty, error) {
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		res, err := union.IsFollow(tx, l.svcCtx.RDS, in.UserId, in.TargetId)
		if err != nil {
			return err
		}
		// res == false means not following yet
		if !res {
			return nil
		}

		err = tx.
			Where("follower_id = ? AND followed_id = ?", in.UserId, in.TargetId).
			Delete(&model.Follow{}).
			Error
		if err != nil {
			return utils.InternalWithDetails("(mysql)error deleting follow relation", err)
		}

		err = l.svcCtx.RDS.Set(cache.GenFollowKey(in.UserId, in.TargetId), "no")
		if err != nil {
			return utils.InternalWithDetails("(redis)error updating follow relation", err)
		}

		// update user counts
		err = union.ModifyUserCounts(tx, l.svcCtx.RDS, in.UserId, "following_count", -1)
		if err != nil {
			return err
		}
		err = union.ModifyUserCounts(tx, l.svcCtx.RDS, in.TargetId, "follower_count", -1)
		if err != nil {
			return err
		}

		// update friend relation
		idA, idB := utils.SortId(in.UserId, in.TargetId)
		err = tx.Where("user_a_id = ? AND user_b_id = ?", idA, idB).Delete(&model.Friend{}).Error
		if err != nil {
			return utils.InternalWithDetails("err deleting friend relation", err)
		}

		return nil
	})
	// transaction end, handle error if occurred
	if err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}
