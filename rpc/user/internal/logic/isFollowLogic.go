package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/internal/union"
	"tikstart/rpc/user/user"
)

type IsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowLogic {
	return &IsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowLogic) IsFollow(in *user.IsFollowRequest) (*user.IsFollowResponse, error) {
	res, err := union.IsFollow(l.svcCtx.DB, l.svcCtx.RDS, in.UserId, in.TargetId)
	if err != nil {
		logx.WithContext(l.ctx).Error(err)
		return nil, err
	}

	return &user.IsFollowResponse{
		IsFollow: res,
	}, nil
}
