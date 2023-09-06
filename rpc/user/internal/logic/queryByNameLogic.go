package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/internal/union"
	"tikstart/rpc/user/user"
)

type QueryByNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryByNameLogic {
	return &QueryByNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryByNameLogic) QueryByName(in *user.QueryByNameRequest) (*user.UserInfo, error) {
	username := in.Username

	userInfo, err := union.GetUserInfoByName(l.svcCtx.DB, l.svcCtx.RDS, username)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
