package logic

import (
	"context"
	"gorm.io/gorm"
	"tikstart/common/model"
	"tikstart/common/utils"

	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModWorkCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModWorkCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModWorkCountLogic {
	return &ModWorkCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModWorkCountLogic) ModWorkCount(in *user.ModWorkCountRequest) (*user.Empty, error) {
	err := l.svcCtx.DB.
		Model(&model.User{}).
		Where("user_id = ?", in.UserId).
		UpdateColumn("work_count", gorm.Expr("work_count + ?", in.Delta)).
		Error

	if err != nil {
		return nil, utils.InternalWithDetails("(mysql)error modifying work_count", err)
	}
	return &user.Empty{}, nil
}
