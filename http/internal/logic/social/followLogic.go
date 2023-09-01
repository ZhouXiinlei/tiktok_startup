package social

import (
	"context"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/user/common"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	Follow   = 1
	UnFollow = 2
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowRequest) (resp *types.FollowResponse, err error) {
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)
	_, err = l.svcCtx.UserRpc.QueryById(l.ctx, &user.QueryByIdRequest{
		UserId: req.ToUserId,
	})
	if _, match := utils.MatchError(err, common.ErrUserNotFound); match {
		return nil, schema.ApiError{
			StatusCode: 422,
			Code:       42202,
			Message:    "关注的用户不存在",
		}
	}
	if err != nil {
		return nil, schema.ServerError{
			ApiError: schema.ApiError{
				StatusCode: 500,
				Code:       50000,
				Message:    "Internal Server Error",
			},
			Detail: err,
		}
	}
	if req.ActionType == Follow {
		res, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:   userClaims.UserId,
			TargetId: req.ToUserId,
		})
		if err != nil {
			return nil, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 500,
					Code:       50000,
					Message:    "Internal Server Error",
				},
				Detail: err,
			}
		}
		isFollow := res.IsFollow
		if isFollow {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42203,
				Message:    "已经关注过了",
			}
		}
		_, err = l.svcCtx.UserRpc.Follow(l.ctx, &user.FollowRequest{
			UserId:   userClaims.UserId,
			TargetId: req.ToUserId,
		})
		if err != nil {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42203,
				Message:    "关注失败",
			}
		}
	}
	if req.ActionType == UnFollow {
		res, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:   userClaims.UserId,
			TargetId: req.ToUserId,
		})
		if err != nil {
			return nil, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 500,
					Code:       50000,
					Message:    "Internal Server Error",
				},
				Detail: err,
			}
		}
		isFollow := res.IsFollow
		if !isFollow {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42203,
				Message:    "没关注过这个用户",
			}
		}
		_, err = l.svcCtx.UserRpc.UnFollow(l.ctx, &user.UnFollowRequest{
			UserId:   userClaims.UserId,
			TargetId: req.ToUserId,
		})
		if err != nil {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42203,
				Message:    "取消关注失败",
			}
		}
	}
	return &types.FollowResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
	}, nil
}
