package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/internal/union"
	"tikstart/rpc/user/user"
)

type QueryByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryByIdLogic {
	return &QueryByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryByIdLogic) QueryById(in *user.QueryByIdRequest) (*user.UserInfo, error) {
	userId := in.UserId

	userInfo, err := union.GetUserInfoById(l.svcCtx.DB, l.svcCtx.RDS, userId)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
