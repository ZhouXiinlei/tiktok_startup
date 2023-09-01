package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"tikstart/common/utils"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/video"
	"tikstart/rpc/video/videoClient"
	"time"

	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentListLogic) GetCommentList(req *types.GetCommentListRequest) (resp *types.GetCommentListResponse, err error) {
	logx.WithContext(l.ctx).Infof("获取评论列表: %+v", req)

	UserClaims, err := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)
	if err != nil {
		return nil, err
	}

	commentListData, err := l.svcCtx.VideoRpc.GetCommentList(l.ctx, &videoClient.GetCommentListRequest{
		VideoId: req.VideoId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取评论列表失败: %v", err)
		return nil, err
	}

	order := make(map[int]int, len(commentListData.CommentList))

	commentList, err := mr.MapReduce(func(source chan<- interface{}) {
		for i, c := range commentListData.CommentList {
			source <- c
			order[int(c.Id)] = i
		}
	}, func(item interface{}, writer mr.Writer[types.Comment], cancel func(error)) {
		comment := item.(*video.Comment)

		userInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
			UserId: comment.AuthorId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户信息失败: %v", err)
			cancel(err)
			return
		}

		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userClient.IsFollowRequest{
			UserId:   UserClaims.UserId,
			TargetId: comment.AuthorId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取关注信息失败: %v", err)
			cancel(err)
			return
		}

		writer.Write(types.Comment{
			Id:         comment.Id,
			Content:    comment.Content,
			CreateDate: time.Unix(comment.CreateTime, 0).Format("01-02"),
			User: types.User{
				Id:            userInfo.UserId,
				Name:          userInfo.Username,
				IsFollow:      isFollowRes.IsFollow,
				FollowCount:   userInfo.FollowCount,
				FollowerCount: userInfo.FollowerCount,
			},
		})
	}, func(pipe <-chan types.Comment, writer mr.Writer[[]types.Comment], cancel func(error)) {
		list := make([]types.Comment, len(commentListData.CommentList))
		for item := range pipe {
			comment := item
			i, ok := order[int(comment.Id)]
			if !ok {
				cancel(err)
				return
			}

			list[i] = comment
		}
		writer.Write(list)
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取评论列表失败: %v", err)
		return nil, err
	}

	return &types.GetCommentListResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
		CommentList: commentList,
	}, nil

}
