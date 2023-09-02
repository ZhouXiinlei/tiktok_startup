package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/common/utils"

	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentByIdLogic {
	return &GetCommentByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentByIdLogic) GetCommentById(in *video.GetCommentByIdRequest) (*video.GetCommentByIdResponse, error) {
	var comment model.Comment
	err := l.svcCtx.Mysql.Where("comment_id = ?", in.CommentId).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrCommentNotFound.Err()
		}
		return nil, utils.InternalWithDetails("err querying comment", err)
	}

	return &video.GetCommentByIdResponse{
		CommentId:   comment.CommentId,
		UserId:      comment.UserId,
		Content:     comment.Content,
		CreatedTime: comment.CreatedAt.Unix(),
	}, nil
}
