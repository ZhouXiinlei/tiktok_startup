package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/svc"
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
	// api should check user existence first, this interface doesn't
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		err := tx.
			Model(&model.Follow{}).
			Where("follower_id = ? AND followed_id = ?", in.UserId, in.TargetId).
			Count(&count).
			Error
		if err != nil {
			return utils.InternalWithDetails("error querying follow record", err)
		}
		// follow record already exists, no need to modify count
		if count > 0 {
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
		return nil
	})

	// transaction end, handle error and return empty
	if err != nil {
		return nil, err
	}
	return &user.Empty{}, nil
}
