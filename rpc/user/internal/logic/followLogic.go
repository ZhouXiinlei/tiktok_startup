package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/cache"
	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/internal/union"
	"tikstart/rpc/user/user"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowLogic) Follow(in *user.FollowRequest) (*user.Empty, error) {
	// api should check User existence first, this interface doesn't
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		res, err := union.IsFollow(tx, l.svcCtx.RDS, in.UserId, in.TargetId)
		if err != nil {
			return err
		}
		// res == true means already following
		if res {
			return nil
		}

		// create follow record
		err = tx.Create(&model.Follow{
			FollowerId: in.UserId,
			FollowedId: in.TargetId,
		}).Error
		if err != nil {
			return utils.InternalWithDetails("(redis)error creating follow record", err)
		}

		err = l.svcCtx.RDS.Set(cache.GenFollowKey(in.UserId, in.TargetId), "yes")
		if err != nil {
			return utils.InternalWithDetails("(redis)error updating follow relation", err)
		}

		return nil
	})
	// transaction end, handle error and return empty
	if err != nil {
		return nil, err
	}

	err = union.ModifyUserCounts(l.svcCtx.DB, l.svcCtx.RDS, in.UserId, "following_count", 1)
	if err != nil {
		return nil, utils.InternalWithDetails("error adding following_count", err)
	}
	err = union.ModifyUserCounts(l.svcCtx.DB, l.svcCtx.RDS, in.TargetId, "follower_count", 1)
	if err != nil {
		return nil, utils.InternalWithDetails("error adding follower_count", err)
	}

	// check friend relation
	res, err := union.IsFollow(l.svcCtx.DB, l.svcCtx.RDS, in.TargetId, in.UserId)
	if err != nil {
		return nil, err
	}
	if res {
		idA, idB := utils.SortId(in.UserId, in.TargetId)
		err = l.svcCtx.DB.Create(&model.Friend{
			UserAId: idA,
			UserBId: idB,
		}).Error
		if err != nil {
			return nil, utils.InternalWithDetails("error creating friend record", err)
		}
	}

	return &user.Empty{}, nil
}
