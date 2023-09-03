package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"tikstart/common"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/video/videoClient"
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
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	switch req.ActionType {
	case Favorite:
		if _, err = l.svcCtx.VideoRpc.FavoriteVideo(l.ctx, &videoClient.FavoriteVideoRequest{
			UserId:  userClaims.UserId,
			VideoId: req.VideoId,
		}); err != nil {
			if st, match := utils.MatchError(err, common.ErrVideoNotFound); match {
				return nil, schema.ApiError{
					StatusCode: 422,
					Code:       42204,
					Message:    "视频不存在",
				}
			} else {
				return nil, utils.ReturnInternalError(st, err)
			}
		}
	case UnFavorite:
		if _, err = l.svcCtx.VideoRpc.UnFavoriteVideo(l.ctx, &videoClient.UnFavoriteVideoRequest{
			UserId:  userClaims.UserId,
			VideoId: req.VideoId,
		}); err != nil {
			st, _ := status.FromError(err)
			return nil, utils.ReturnInternalError(st, err)
		}
	default:
		return nil, schema.ApiError{
			StatusCode: 422,
			Code:       42205,
			Message:    "未知操作",
		}
	}

	return &types.FavoriteResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
	}, nil
}
