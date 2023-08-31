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
		res := tx.Clauses(clause.OnConflict{
			// mysql not supports "on conflict(primary key)" feature,
			// so we must remove follow_id column from table,
			// and set primary key to follower_id and followed_id,
			// then it will check if primary key is duplicated,
			// so the clause can function normally.
			// the next line is useless , but I decide to keep it.
			Columns: []clause.Column{{Name: "follower_id"}, {Name: "followed_id"}},
			// not updating time to make follow idempotent
			//DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
			DoNothing: true,
		}).Create(&model.Follow{
			FollowerId: in.UserId,
			FollowedId: in.TargetId,
		})

		// firstly check error
		if err := res.Error; err != nil {
			return utils.InternalWithDetails("error creating follow relation", err)
		}
		// if no rows affected, meaning follow relation already exists
		if res.RowsAffected == 0 {
			return nil
		}

		// codes going here, meaning no error occurred, should add count
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.User{}).Where("user_id = ?", in.UserId).UpdateColumn("following_count", gorm.Expr("following_count + ?", 1)).Error
		if err != nil {
			return utils.InternalWithDetails("error adding following_count", err)
		}

		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.User{}).Where("user_id = ?", in.TargetId).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
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
