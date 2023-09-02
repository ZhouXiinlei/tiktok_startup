package logic

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tikstart/common/model"
	"tikstart/common/utils"

	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *video.CreateCommentRequest) (*video.CreateCommentResponse, error) {
	comment := model.Comment{
		VideoId: in.VideoId,
		UserId:  in.UserId,
		Content: in.Content,
	}

	err := l.svcCtx.Mysql.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&comment).Error; err != nil {
			return utils.InternalWithDetails("error creating comment", err)
		}

		if err := tx.Model(&model.Video{}).
			Where("video_id = ?", in.VideoId).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).
			Error; err != nil {
			return utils.InternalWithDetails("error updating comment count", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &video.CreateCommentResponse{
		CommentId:   comment.CommentId,
		UserId:      comment.UserId,
		Content:     comment.Content,
		CreatedTime: comment.CreatedAt.Unix(),
	}, nil
}
