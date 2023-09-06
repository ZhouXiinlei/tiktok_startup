package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
	"tikstart/common"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/videoClient"
	"time"
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
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	switch req.ActionType {
	case Publish:
		regPattern := regexp.MustCompile("\\s+")
		commentText := regPattern.ReplaceAllString(req.CommentText, "")
		commentText = strings.Trim(commentText, " ")

		if commentText == "" {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42208,
				Message:    "评论内容不能为空",
			}
		}

		var comment types.Comment

		err := mr.Finish(func() (err error) {
			createResp, err := l.svcCtx.VideoRpc.CreateComment(l.ctx, &videoClient.CreateCommentRequest{
				UserId:  userClaims.UserId,
				VideoId: req.VideoId,
				Content: commentText,
			})
			if err != nil {
				if st, match := utils.MatchError(err, common.ErrVideoNotFound); match {
					return schema.ApiError{
						StatusCode: 422,
						Code:       42209,
						Message:    "视频不存在",
					}
				} else {
					logx.WithContext(l.ctx).Errorf("创建评论失败: %s", err.Error())
					return utils.ReturnInternalError(l.ctx, st, err)
				}
			}

			comment.Content = createResp.Content
			comment.Id = createResp.CommentId
			comment.CreateDate = time.Unix(createResp.CreatedTime, 0).Format("01-02")

			return nil
		}, func() error {
			userInfo, err := l.svcCtx.UserRpc.QueryById(l.ctx, &userClient.QueryByIdRequest{
				UserId: userClaims.UserId,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("获取用户信息失败: %s", err.Error())
				return utils.ReturnInternalError(l.ctx, status.Convert(err), err)
			}

			comment.User = types.User{
				Id:             userInfo.UserId,
				Name:           userInfo.Username,
				IsFollow:       false,
				FollowCount:    userInfo.FollowingCount,
				FollowerCount:  userInfo.FollowerCount,
				TotalFavorited: userInfo.TotalFavorited,
				WorkCount:      userInfo.WorkCount,
				FavoriteCount:  userInfo.FavoriteCount,
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		return &types.CommentResponse{
			BasicResponse: types.BasicResponse{
				StatusCode: 0,
				StatusMsg:  "Success",
			},
			Comment: comment,
		}, nil
	case Delete:
		commentInfo, err := l.svcCtx.VideoRpc.GetCommentById(l.ctx, &videoClient.GetCommentByIdRequest{
			CommentId: req.CommentId,
		})
		if err != nil {
			if st, match := utils.MatchError(err, common.ErrCommentNotFound); match {
				return nil, schema.ApiError{
					StatusCode: 422,
					Code:       42210,
					Message:    "评论不存在",
				}
			} else {
				logx.WithContext(l.ctx).Errorf("获取评论信息失败: %s", err.Error())
				return nil, utils.ReturnInternalError(l.ctx, st, err)
			}
		}

		if commentInfo.UserId != userClaims.UserId {
			return nil, schema.ApiError{
				StatusCode: 403,
				Code:       40300,
				Message:    "用户无权限删除此评论",
			}
		}

		if _, err = l.svcCtx.VideoRpc.DeleteComment(l.ctx, &videoClient.DeleteCommentRequest{
			CommentId: req.CommentId,
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("删除评论失败: %s", err.Error())
			return nil, utils.ReturnInternalError(l.ctx, status.Convert(err), err)
		}

		return &types.CommentResponse{
			BasicResponse: types.BasicResponse{
				StatusCode: 0,
				StatusMsg:  "Success",
			},
		}, nil
	default:
		return nil, schema.ApiError{
			StatusCode: 422,
			Code:       42200,
			Message:    "未知操作",
		}
	}

}
