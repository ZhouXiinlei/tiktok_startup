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
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		follower := model.User{}
		followed := model.User{}

		// api should check user existence first
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", in.UserId).First(&follower)
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", in.WillFollow).First(&followed)

		follower.FollowingCount++
		err := tx.Save(&follower).Error
		if err != nil {
			return utils.InternalWithDetails("error saving following number", err)
		}

		followed.FollowerCount++
		err = tx.Save(&followed).Error
		if err != nil {
			return utils.InternalWithDetails("error saving follower number", err)
		}

		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&follower).Association("Following").Replace(&followed)
		if err != nil {
			return utils.InternalWithDetails("error adding follow relation", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}
