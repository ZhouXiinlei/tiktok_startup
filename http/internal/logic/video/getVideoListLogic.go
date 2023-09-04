package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/videoClient"
	"time"
)

type GetVideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListLogic {
	return &GetVideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVideoListLogic) GetVideoList(req *types.GetVideoListRequest) (resp *types.GetVideoListResponse, err error) {
	userClaims, err := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)
	var userId int64 = 0
	if err == nil {
		userId = userClaims.UserId
	}

	latestTime := time.Now().Unix()
	if req.LatestTime != 0 {
		latestTime = req.LatestTime / 1000
	}

	videoListRes, err := l.svcCtx.VideoRpc.GetVideoList(l.ctx, &videoClient.GetVideoListRequest{
		Num:        20,
		LatestTime: latestTime,
	})
	if err != nil {
		return nil, utils.ReturnInternalError(status.Convert(err), err)
	}

	// TODO: 直接让RPC返回NextTime
	var nextTime int64 = 0
	if len(videoListRes.VideoList) != 0 {
		nextTime = videoListRes.VideoList[len(videoListRes.VideoList)-1].CreateTime
	}

	// 补充视频信息
	videoList := make([]types.Video, 0, len(videoListRes.VideoList))
	for _, v := range videoListRes.VideoList {
		//获取作者信息和关注情况
		userInfoResp, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
			UserId: v.AuthorId,
		})
		if err != nil {
			return nil, utils.ReturnInternalError(status.Convert(err), err)
		}
		//获取视频收藏状态
		isFavorite := false
		if userId != 0 {
			isFavoriteVideoRes, err := l.svcCtx.VideoRpc.IsFavoriteVideo(l.ctx, &videoClient.IsFavoriteVideoRequest{
				UserId:  userId,
				VideoId: v.Id,
			})
			if err != nil {
				return nil, utils.ReturnInternalError(status.Convert(err), err)
			}
			isFavorite = isFavoriteVideoRes.IsFavorite
		}
		//获取作者关注状态
		isFollow := false
		if userId != 0 {
			isFollowVideoRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userClient.IsFollowRequest{
				UserId:   userId,
				TargetId: v.AuthorId,
			})
			if err != nil {
				return nil, utils.ReturnInternalError(status.Convert(err), err)
			}
			isFollow = isFollowVideoRes.IsFollow

		}
		videoList = append(videoList, types.Video{
			Id:    v.Id,
			Title: v.Title,
			Author: types.User{
				Id:            v.AuthorId,
				Name:          userInfoResp.Username,
				IsFollow:      isFollow,
				FollowCount:   userInfoResp.FollowingCount,
				FollowerCount: userInfoResp.FollowerCount,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavorite,
		})
	}

	return &types.GetVideoListResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
		VideoList: videoList,
		NextTime:  nextTime,
	}, nil
}
