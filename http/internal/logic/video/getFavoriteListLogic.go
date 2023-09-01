package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"tikstart/common/utils"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/video"

	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/rpc/video/videoClient"

	"github.com/zeromicro/go-zero/core/logx"
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
	UserClaims, err := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)
	if err != nil {
		return nil, err
	}

	GetFavoriteListResponse, err := l.svcCtx.VideoRpc.GetFavoriteVideoList(l.ctx, &videoClient.GetFavoriteVideoListRequest{
		UserId: UserClaims.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户喜欢视频列表失败: %v", err)
		return nil, err
	}

	order := make(map[int]int, len(GetFavoriteListResponse.VideoList))

	favoriteList, err := mr.MapReduce(func(source chan<- interface{}) {
		for i, v := range GetFavoriteListResponse.VideoList {
			source <- v
			order[int(v.Id)] = i
		}
	}, func(item interface{}, writer mr.Writer[types.Video], cancel func(error)) {
		videoInfo := item.(*video.VideoInfo)

		authorInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
			UserId: videoInfo.AuthorId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取作者信息失败: %v", err)
			cancel(err)
		}

		isFollow, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userClient.IsFollowRequest{
			UserId:   UserClaims.UserId,
			TargetId: authorInfo.UserId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取关注信息失败: %v", err)
			cancel(err)
		}

		author := types.User{
			Id:            authorInfo.UserId,
			Name:          authorInfo.Username,
			FollowCount:   authorInfo.FollowCount,
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
		list := make([]types.Video, len(GetFavoriteListResponse.VideoList))
		for item := range pipe {
			videoInfo := item
			i, err := order[int(videoInfo.Id)]
			if !err {
				return
			}
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
