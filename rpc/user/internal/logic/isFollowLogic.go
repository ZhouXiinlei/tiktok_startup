package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/cache"
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
	// cache miss with return ""
	val, err := l.svcCtx.RDS.Get(cache.GenFollowKey(in.UserId, in.TargetId))
	if err != nil {
		return nil, utils.InternalWithDetails("(redis)error getting follow relation", err)
	}
	if val != "" {
		if val == "yes" {
			return &user.IsFollowResponse{
				IsFollow: true,
			}, nil
		} else {
			return &user.IsFollowResponse{
				IsFollow: false,
			}, nil
		}
	}

	var count int64
	err = l.svcCtx.DB.
		Model(&model.Follow{}).
		Where("follower_id = ? AND followed_id = ?", in.UserId, in.TargetId).
		Count(&count).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("(mysql)error querying follow relation", err)
	}

	err = l.svcCtx.RDS.Set(cache.GenFollowKey(in.UserId, in.TargetId), cache.YesOrNo(count == 1))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("(redis)error setting follow relation", err)
	}

	return &user.IsFollowResponse{
		IsFollow: count == 1,
	}, nil
}
