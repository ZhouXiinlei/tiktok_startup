package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"tikstart/common"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/user/user"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoRequest) (resp *types.GetUserInfoResponse, err error) {
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)
	targetId := req.UserId

	userInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &user.QueryByIdRequest{
		UserId: targetId,
	})
	if err != nil {
		if st, match := utils.MatchError(err, common.ErrUserNotFound); match {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42202,
				Message:    "用户名不存在",
			}
		} else {
			return nil, utils.ReturnInternalError(l.ctx, st, err)
		}
	}

	isFollowResp, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
		UserId:   userClaims.UserId,
		TargetId: targetId,
	})
	if err != nil {
		return nil, utils.ReturnInternalError(l.ctx, status.Convert(err), err)
	}

	return &types.GetUserInfoResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		User: types.User{
			Id:             userInfo.UserId,
			Name:           userInfo.Username,
			FollowCount:    userInfo.FollowingCount,
			FollowerCount:  userInfo.FollowerCount,
			IsFollow:       isFollowResp.IsFollow,
			TotalFavorited: userInfo.TotalFavorited,
			WorkCount:      userInfo.WorkCount,
			FavoriteCount:  userInfo.FavoriteCount,
		},
	}, nil
}
