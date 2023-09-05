package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"
)

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentListLogic) GetCommentList(in *video.GetCommentListRequest) (*video.GetCommentListResponse, error) {
	var count int64
	err := l.svcCtx.DB.Model(&model.Video{}).Where("video_id = ?", in.VideoId).Count(&count).Error
	if err != nil {
		return nil, utils.InternalWithDetails("err querying video info", err)
	}
	if count == 0 {
		return nil, common.ErrVideoNotFound.Err()
	}

	var comments []*model.Comment
	if err := l.svcCtx.DB.
		Where("video_id = ?", in.VideoId).
		Preload("User").
		Order("created_at").
		Find(&comments).Error; err != nil {
		return nil, utils.InternalWithDetails("err querying comment list", err)
	}

	commentList := make([]*video.Comment, 0, len(comments))
	for _, comment := range comments {
		commentList = append(commentList, &video.Comment{
			Id:             comment.CommentId,
			AuthorId:       comment.UserId,
			CreateTime:     comment.CreatedAt.Unix(),
			Content:        comment.Content,
			UserId:         comment.User.UserId,
			Username:       comment.User.Username,
			FollowingCount: comment.User.FollowingCount,
			FollowerCount:  comment.User.FollowerCount,
		})
	}

	return &video.GetCommentListResponse{
		CommentList: commentList,
	}, nil
}
