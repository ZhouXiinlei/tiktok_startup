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

func (l *QueryByIdLogic) QueryById(in *user.QueryByIdRequest) (*user.QueryResponse, error) {
	userId := in.UserId

	userRecord := model.User{}
	err := l.svcCtx.DB.Where("user_id = ?", userId).First(&userRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound.Err()
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	followingCount, err := union.PickUserCounts(l.svcCtx.DB, l.svcCtx.RDS, in.UserId, "following_count", userRecord.FollowingCount)
	if err != nil {
		return nil, err
	}
	followerCount, err := union.PickUserCounts(l.svcCtx.DB, l.svcCtx.RDS, in.UserId, "follower_count", userRecord.FollowerCount)
	if err != nil {
		return nil, err
	}

	return &user.QueryResponse{
		UserId:         userRecord.UserId,
		Username:       userRecord.Username,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
		Password:       userRecord.Password,
		CreatedAt:      userRecord.CreatedAt.Unix(),
		UpdatedAt:      userRecord.UpdatedAt.Unix(),
	}, nil
}
