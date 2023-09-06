package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"tikstart/common"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/user/user"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	username := req.Username
	password := req.Password

	res, err := l.svcCtx.UserRpc.QueryByName(l.ctx, &user.QueryByNameRequest{
		Username: username,
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

	err = bcrypt.CompareHashAndPassword(res.Password, []byte(password))
	if err != nil {
		return nil, schema.ApiError{
			StatusCode: 422,
			Code:       42203,
			Message:    "密码错误",
		}
	}

	tokenString, err := utils.CreateToken(res.UserId, l.svcCtx.Config.JwtAuth.Secret, l.svcCtx.Config.JwtAuth.Expire)
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

	return &types.LoginResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserId: res.UserId,
		Token:  tokenString,
	}, nil
}
