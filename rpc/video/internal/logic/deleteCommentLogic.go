package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/cache"

	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *video.DeleteCommentRequest) (*video.Empty, error) {
	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		// MySQL doesn't support returning feature, so we must select the comment first
		err := tx.
			Where("comment_id = ?", in.CommentId).
			First(&comment).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.ErrCommentNotFound.Err()
			}
			return utils.InternalWithDetails("error querying comment", err)
		}

		if err := tx.Where("comment_id = ?", in.CommentId).Delete(&comment).Error; err != nil {
			return utils.InternalWithDetails("error deleting comment", err)
		}

		err = cache.ModifyVideoCounts(tx, l.svcCtx.RDS, in.CommentId, "comment_count", -1)
		if err != nil {
			return utils.InternalWithDetails("error deleting comment_count", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &video.Empty{}, nil
}
