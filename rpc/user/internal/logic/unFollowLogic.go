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
		follower := model.User{}
		unfollowed := model.User{}

		// api should check user existence first
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", in.UserId).First(&follower)
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", in.WillUnfollow).First(&unfollowed)

		follower.FollowingCount--
		err := tx.Save(&follower).Error
		if err != nil {
			return utils.InternalWithDetails("error saving following number", err)
		}

		unfollowed.FollowerCount--
		err = tx.Save(&unfollowed).Error
		if err != nil {
			return utils.InternalWithDetails("error saving follower number", err)
		}

		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&follower).Association("Following").Delete(&unfollowed)
		if err != nil {
			return utils.InternalWithDetails("error deleting follow relation", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}
