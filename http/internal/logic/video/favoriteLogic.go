package video

import (
	"context"

	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/rpc/video/videoClient"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	Favorite   = 1
	UnFavorite = 2
)

type FavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteLogic) Favorite(req *types.FavoriteRequest) (resp *types.FavoriteResponse, err error) {
	logx.WithContext(l.ctx).Infof("收藏视频: %v", req)

	userClaims, err := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)
	if err != nil {
		return nil, err
	}

	if req.ActionType == Favorite {
		if _, err = l.svcCtx.VideoRpc.FavoriteVideo(l.ctx, &videoClient.FavoriteVideoRequest{
			UserId:  userClaims.UserId,
			VideoId: req.VideoId,
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("收藏视频失败: %v", err)
			return nil, err
		}

	} else if req.ActionType == UnFavorite {
		if _, err = l.svcCtx.VideoRpc.UnFavoriteVideo(l.ctx, &videoClient.UnFavoriteVideoRequest{
			UserId:  userClaims.UserId,
			VideoId: req.VideoId,
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("收藏视频失败: %v", err)
			return nil, err
		}

	} else {
		return nil, err
	}

	return &types.FavoriteResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
	}, nil

}
