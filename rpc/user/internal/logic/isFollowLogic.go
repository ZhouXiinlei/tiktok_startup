package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/common/utils"

	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
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
	var count int64
	err := l.svcCtx.DB.Model(&model.Follow{}).Where("follower_id = ? AND followed_id = ?", in.UserId, in.TargetId).Count(&count).Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying follow relation", err)
	}

	return &user.IsFollowResponse{
		IsFollow: count == 1,
	}, nil
}
