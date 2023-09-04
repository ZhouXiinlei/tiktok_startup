package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"google.golang.org/grpc/status"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/video"
	"tikstart/rpc/video/videoClient"
)

type GetFavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFavoriteListLogic) GetFavoriteList(req *types.GetFavoriteListRequest) (resp *types.GetFavoriteListResponse, err error) {
	logx.WithContext(l.ctx).Infof("获取用户喜欢视频列表: %+v", req)

	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	favoriteListRes, err := l.svcCtx.VideoRpc.GetFavoriteVideoList(l.ctx, &videoclient.GetFavoriteVideoListRequest{
		UserId: userClaims.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户喜欢视频列表失败: %v", err)
		return nil, utils.ReturnInternalError(status.Convert(err), err)
	}

	order := make(map[int]int, len(favoriteListRes.VideoList))

	favoriteList, err := mr.MapReduce(func(source chan<- *video.VideoInfo) {
		for i, v := range favoriteListRes.VideoList {
			source <- v
			order[int(v.Id)] = i
		}
	}, func(videoInfo *video.VideoInfo, writer mr.Writer[types.Video], cancel func(error)) {
		authorInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
			UserId: videoInfo.AuthorId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取作者信息失败: %v", err)
			cancel(utils.ReturnInternalError(status.Convert(err), err))
			return
		}

		isFollow, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userClient.IsFollowRequest{
			UserId:   userClaims.UserId,
			TargetId: authorInfo.UserId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取关注信息失败: %v", err)
			cancel(utils.ReturnInternalError(status.Convert(err), err))
			return
		}

		// TODO: RPC使用PRELOAD直接查出作者信息
		author := types.User{
			Id:            authorInfo.UserId,
			Name:          authorInfo.Username,
			FollowCount:   authorInfo.FollowingCount,
			FollowerCount: authorInfo.FollowerCount,
			IsFollow:      isFollow.IsFollow,
		}

		writer.Write(types.Video{
			Id:            videoInfo.Id,
			Title:         videoInfo.Title,
			Author:        author,
			PlayUrl:       videoInfo.PlayUrl,
			CoverUrl:      videoInfo.CoverUrl,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    true,
		})

	}, func(pipe <-chan types.Video, writer mr.Writer[[]types.Video], cancel func(error)) {
		list := make([]types.Video, len(favoriteListRes.VideoList))
		for item := range pipe {
			videoInfo := item
			i, _ := order[int(videoInfo.Id)]
			list[i] = videoInfo
		}
		writer.Write(list)
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户喜欢视频列表失败: %v", err)
		return nil, err
	}

	return &types.GetFavoriteListResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
		VideoList: favoriteList,
	}, nil
}
