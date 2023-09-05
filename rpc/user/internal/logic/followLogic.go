package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tikstart/common/model"
	"tikstart/common/utils"
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
		res, err := union.IsFollow(l.svcCtx, in.UserId, in.TargetId)
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
			return utils.InternalWithDetails("error creating follow record", err)
		}

		// modify count
		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&model.User{}).
			Where("user_id = ?", in.UserId).
			UpdateColumn("following_count", gorm.Expr("following_count + ?", 1)).
			Error
		if err != nil {
			return utils.InternalWithDetails("error adding following_count", err)
		}

		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.User{}).
			Where("user_id = ?", in.TargetId).
			UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).
			Error
		if err != nil {
			return utils.InternalWithDetails("error adding follower_count", err)
		}

		// check friend relation
		var count int64
		err = l.svcCtx.DB.Model(&model.Follow{}).Where("follower_id = ? AND followed_id = ?", in.TargetId, in.UserId).Count(&count).Error
		if err != nil {
			return utils.InternalWithDetails("error querying friend record", err)
		}
		if count > 0 {
			idA, idB := utils.SortId(in.UserId, in.TargetId)
			err = tx.Create(&model.Friend{
				UserAId: idA,
				UserBId: idB,
			}).Error
			if err != nil {
				return utils.InternalWithDetails("error creating friend record", err)
			}
		}
		return nil
	})

	// transaction end, handle error and return empty
	if err != nil {
		return nil, err
	}
	return &user.Empty{}, nil
}
