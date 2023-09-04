package logic

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tikstart/common/model"
	"tikstart/common/utils"

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
		res := tx.Where("follower_id = ? AND followed_id = ?", in.UserId, in.TargetId).Delete(&model.Follow{})
		if err := res.Error; err != nil {
			return utils.InternalWithDetails("error deleting follow record", err)
		}
		if res.RowsAffected == 0 {
			return nil
		}

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&model.User{}).
			Where("user_id = ?", in.UserId).
			UpdateColumn("following_count", gorm.Expr("following_count - ?", 1)).
			Error
		if err != nil {
			return utils.InternalWithDetails("error reducing following_count", err)
		}

		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&model.User{}).
			Where("user_id = ?", in.TargetId).
			UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).
			Error
		if err != nil {
			return utils.InternalWithDetails("error reducing follower_count", err)
		}

		// update friend relation
		idA, idB := utils.SortId(in.UserId, in.TargetId)
		err = l.svcCtx.DB.Where("user_a_id = ? AND user_b_id = ?", idA, idB).Delete(&model.Friend{}).Error
		if err != nil {
			return utils.InternalWithDetails("err deleting friend relation", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &user.Empty{}, nil
}
