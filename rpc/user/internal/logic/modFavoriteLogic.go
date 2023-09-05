package logic

import (
	"context"
	"tikstart/rpc/user/internal/union"

	"tikstart/rpc/user/internal/svc"
	"tikstart/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModFavoriteLogic {
	return &ModFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModFavoriteLogic) ModFavorite(in *user.ModFavoriteRequest) (*user.Empty, error) {

	err := union.ModifyUserCounts(l.svcCtx.DB, l.svcCtx.RDS, in.UserId, "favorite_count", in.Delta)
	if err != nil {
		return nil, err
	}

	err = union.ModifyUserCounts(l.svcCtx.DB, l.svcCtx.RDS, in.TargetId, "total_favorited", in.Delta)
	if err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}
