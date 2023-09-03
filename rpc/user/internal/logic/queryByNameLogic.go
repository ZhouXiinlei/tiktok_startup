package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/rpc/user/internal/svc"
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

func (l *QueryByNameLogic) QueryByName(in *user.QueryByNameRequest) (*user.QueryResponse, error) {
	username := in.Username

	userRecord := model.User{}
	err := l.svcCtx.DB.Where("username = ?", username).First(&userRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound.Err()
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &user.QueryResponse{
		UserId:         userRecord.UserId,
		Username:       userRecord.Username,
		FollowingCount: userRecord.FollowingCount,
		FollowerCount:  userRecord.FollowerCount,
		CreatedAt:      userRecord.CreatedAt.Unix(),
		UpdatedAt:      userRecord.UpdatedAt.Unix(),
	}, nil
}
