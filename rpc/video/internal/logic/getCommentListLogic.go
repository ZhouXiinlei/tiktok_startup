package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
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
	var comments []*model.Comment
	if err := l.svcCtx.DB.
		Where("video_id = ?", in.VideoId).
		Order("created_at").
		Find(&comments).Error; err != nil {
		return nil, utils.InternalWithDetails("err querying comment list", err)
	}

	commentList := make([]*video.Comment, 0, len(comments))
	for _, comment := range comments {
		commentList = append(commentList, &video.Comment{
			Id:         comment.CommentId,
			AuthorId:   comment.UserId,
			CreateTime: comment.CreatedAt.Unix(),
			Content:    comment.Content,
		})
	}

	return &video.GetCommentListResponse{
		CommentList: commentList,
	}, nil
}
