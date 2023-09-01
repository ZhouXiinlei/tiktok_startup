package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"tikstart/common/utils"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/videoClient"
	"time"

	//"tikstart/common/utils"
	//"tikstart/rpc/user/user"
	//"tikstart/rpc/user/userClient"
	//"tikstart/rpc/video/videoClient"

	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	Publish = 1
	Delete  = 2
)

type CommentVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentVideoLogic {
	return &CommentVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentVideoLogic) CommentVideo(req *types.CommentRequest) (resp *types.CommentResponse, err error) {
	logx.WithContext(l.ctx).Infof("评论视频: %v", req)

	userClaims, err := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)
	if err != nil {
		return nil, err
	}

	if req.ActionType == Publish {
		var Comment types.Comment
		err := mr.Finish(func() (err error) {
			res, err := l.svcCtx.VideoRpc.CommentVideo(l.ctx, &videoClient.CommentVideoRequest{
				UserId:  userClaims.UserId,
				VideoId: req.VideoId,
				Content: req.CommentText,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建评论失败: %s", err.Error())
				return
			}
			Comment.Content = res.Content
			Comment.Id = res.Id
			Comment.CreateDate = time.Unix(res.CreatedTime, 0).Format("01-02")
			return nil
		}, func() error {
			userInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
				UserId: userClaims.UserId,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("获取用户信息失败: %s", err.Error())
				return err
			}
			Comment.User = types.User{
				Id:            userInfo.UserId,
				Name:          userInfo.Username,
				IsFollow:      false,
				FollowCount:   userInfo.FollowCount,
				FollowerCount: userInfo.FollowerCount,
			}
			return nil
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("创建评论失败: %s", err.Error())
			return nil, err
		}

		return &types.CommentResponse{
			BasicResponse: types.BasicResponse{
				StatusCode: 0,
				StatusMsg:  "Success",
			},
			Comment: Comment,
		}, nil
	} else if req.ActionType == Delete {
		commentInfo, err := l.svcCtx.VideoRpc.GetCommentInfo(l.ctx, &videoClient.GetCommentInfoRequest{
			CommentId: req.CommentId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取评论信息失败: %s", err.Error())
			return nil, err
		}

		if commentInfo.UserId != userClaims.UserId {
			logx.WithContext(l.ctx).Errorf("用户无权限删除此评论")
			return nil, err
		}

		if _, err = l.svcCtx.VideoRpc.DeleteVideoComment(l.ctx, &videoClient.DeleteVideoCommentRequest{
			CommentId: req.CommentId,
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("删除评论失败: %s", err.Error())
			return nil, err
		}

		return &types.CommentResponse{
			BasicResponse: types.BasicResponse{
				StatusCode: 0,
				StatusMsg:  "Success",
			},
		}, nil
	}
	return nil, err
}
