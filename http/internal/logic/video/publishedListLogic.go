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
)

type PublishedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishedListLogic {
	return &PublishedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishedListLogic) PublishedList(req *types.PublishedListRequest) (resp *types.PublishedListResponse, err error) {
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	publishedList, err := l.svcCtx.VideoRpc.GetVideoListByAuthor(l.ctx, &videoClient.GetVideoListByAuthorRequest{
		AuthorId: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取发布列表失败: %v", err)
		return nil, utils.ReturnInternalError(status.Convert(err), err)
	}

	authorInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户信息失败: %v", err)
		return nil, utils.ReturnInternalError(status.Convert(err), err)
	}

	videoList := make([]types.Video, 0, len(publishedList.Video))

	//是否关注目标用户
	isFollow := false
	if userClaims.UserId != req.UserId {
		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userClient.IsFollowRequest{
			UserId:   userClaims.UserId,
			TargetId: req.UserId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户是否关注失败: %v", err)
			return nil, utils.ReturnInternalError(status.Convert(err), err)
		}
		isFollow = isFollowRes.IsFollow
	}

	for _, v := range publishedList.Video {
		//是否点赞，作者可以给自己点赞
		isFavoriteResponse, err := l.svcCtx.VideoRpc.IsFavoriteVideo(l.ctx, &videoClient.IsFavoriteVideoRequest{
			UserId:  userClaims.UserId,
			VideoId: v.Id,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户是否点赞失败: %v", err)
			return nil, utils.ReturnInternalError(status.Convert(err), err)
		}

		videoList = append(videoList, types.Video{
			Id:    v.Id,
			Title: v.Title,
			Author: types.User{
				Id:             authorInfo.UserId,
				Name:           authorInfo.Username,
				FollowerCount:  authorInfo.FollowerCount,
				FollowCount:    authorInfo.FollowingCount,
				IsFollow:       isFollow,
				TotalFavorited: authorInfo.TotalFavorited,
				FavoriteCount:  authorInfo.FavoriteCount,
				WorkCount:      authorInfo.WorkCount,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavoriteResponse.IsFavorite,
		})
	}

	return &types.PublishedListResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
		VideoList: videoList,
	}, nil
}
