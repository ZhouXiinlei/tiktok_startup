package video

import (
	"context"
	"tikstart/common/utils"
	"tikstart/http/schema"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line
	Userclaims, err := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	publishedList, err := l.svcCtx.VideoRpc.GetVideoListByAuthor(l.ctx, &videoClient.GetVideoListByAuthorRequest{
		AuthorId: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取发布列表失败: %v", err)
		return nil, schema.ServerError{
			ApiError: schema.ApiError{
				StatusCode: 500,
				Code:       50000,
				Message:    "Internal Server Error",
			},
			Detail: err,
		}
	}
	authorInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户信息失败: %v", err)
		return nil, schema.ServerError{
			ApiError: schema.ApiError{
				StatusCode: 500,
				Code:       50000,
				Message:    "Internal Server Error",
			},
			Detail: err,
		}
	}
	videoList := make([]types.Video, 0, len(publishedList.Video))
	if Userclaims.UserId == req.UserId {
		for _, v := range publishedList.Video {
			isFavorite := false
			IsFavoriteVideoResponse, err := l.svcCtx.VideoRpc.IsFavoriteVideo(l.ctx, &videoClient.IsFavoriteVideoRequest{
				UserId:  Userclaims.UserId,
				VideoId: v.Id,
			})
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
			isFavorite = IsFavoriteVideoResponse.IsFavorite

			videoList = append(videoList, types.Video{
				Id:    v.Id,
				Title: v.Title,
				Author: types.User{
					Id:            authorInfo.UserId,
					Name:          authorInfo.Username,
					FollowerCount: authorInfo.FollowerCount,
					FollowCount:   authorInfo.FollowCount,
					IsFollow:      false,
				},
				PlayUrl:       v.PlayUrl,
				CoverUrl:      v.CoverUrl,
				FavoriteCount: v.FavoriteCount,
				CommentCount:  v.CommentCount,
				IsFavorite:    isFavorite,
			})
		}
	} else {
		isFollowResponse, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userClient.IsFollowRequest{
			UserId:   Userclaims.UserId,
			TargetId: req.UserId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户是否关注失败: %v", err)
			return nil, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 500,
					Code:       50000,
					Message:    "Internal Server Error",
				},
				Detail: err,
			}
		}
		for _, v := range publishedList.Video {
			//是否点赞
			isFavoriteResponse, err := l.svcCtx.VideoRpc.IsFavoriteVideo(l.ctx, &videoClient.IsFavoriteVideoRequest{
				UserId:  Userclaims.UserId,
				VideoId: v.Id,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("获取用户是否点赞失败: %v", err)
				return nil, schema.ServerError{
					ApiError: schema.ApiError{
						StatusCode: 500,
						Code:       50000,
						Message:    "Internal Server Error",
					},
					Detail: err,
				}
			}
			videoList = append(videoList, types.Video{
				Id:    v.Id,
				Title: v.Title,
				Author: types.User{
					Id:            authorInfo.UserId,
					Name:          authorInfo.Username,
					FollowerCount: authorInfo.FollowerCount,
					FollowCount:   authorInfo.FollowCount,
					IsFollow:      isFollowResponse.IsFollow,
				},
				PlayUrl:       v.PlayUrl,
				CoverUrl:      v.CoverUrl,
				FavoriteCount: v.FavoriteCount,
				CommentCount:  v.CommentCount,
				IsFavorite:    isFavoriteResponse.IsFavorite,
			})
		}
	}
	return &types.PublishedListResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
		VideoList: videoList,
	}, nil
}
