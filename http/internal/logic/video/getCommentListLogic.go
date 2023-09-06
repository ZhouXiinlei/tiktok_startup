package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"google.golang.org/grpc/status"
	"tikstart/common"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/video"
	"tikstart/rpc/video/videoClient"
	"time"
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

	UserClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	commentListRes, err := l.svcCtx.VideoRpc.GetCommentList(l.ctx, &videoClient.GetCommentListRequest{
		VideoId: req.VideoId,
	})
	if err != nil {
		if st, match := utils.MatchError(err, common.ErrVideoNotFound); match {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42211,
				Message:    "视频不存在",
			}
		} else {
			return nil, utils.ReturnInternalError(l.ctx, st, err)
		}
	}

	//var order sync.Map
	order := make(map[int64]int, len(commentListRes.CommentList))
	lock := make(chan struct{})
	commentList, err := mr.MapReduce(func(source chan<- interface{}) {
		for i, c := range commentListRes.CommentList {
			source <- c
			order[c.Id] = i
			//order.Store(c.Id, i)
		}
		lock <- struct{}{}
	}, func(item interface{}, writer mr.Writer[types.Comment], cancel func(error)) {
		comment := item.(*video.Comment)

		//userInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
		//	UserId: comment.AuthorId,
		//})
		//if err != nil {
		//	logx.WithContext(l.ctx).Errorf("获取用户信息失败: %v", err)
		//	cancel(utils.ReturnInternalError(status.Convert(err), err))
		//	return
		//}

		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userClient.IsFollowRequest{
			UserId:   UserClaims.UserId,
			TargetId: comment.AuthorId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取关注信息失败: %v", err)
			cancel(utils.ReturnInternalError(l.ctx, status.Convert(err), err))
			return
		}
		writer.Write(types.Comment{
			Id:         comment.Id,
			Content:    comment.Content,
			CreateDate: time.Unix(comment.CreateTime, 0).Format("01-02"),
			User: types.User{
				Id:             comment.UserId,
				Name:           comment.Username,
				IsFollow:       isFollowRes.IsFollow,
				FollowCount:    comment.FollowingCount,
				FollowerCount:  comment.FollowerCount,
				FavoriteCount:  comment.FavoriteCount,
				TotalFavorited: comment.Total_Favorited,
				WorkCount:      comment.WorkCount,
			},
		})
	}, func(pipe <-chan types.Comment, writer mr.Writer[[]types.Comment], cancel func(error)) {
		<-lock
		list := make([]types.Comment, len(commentListRes.CommentList))
		for item := range pipe {
			comment := item
			i, _ := order[comment.Id]
			list[i] = comment
			//i, _ := order.Load(comment.Id)
			//list[i.(int)] = comment
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
