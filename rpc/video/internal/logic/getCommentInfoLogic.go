package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/video"
)

type GetCommentInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentInfoLogic {
	return &GetCommentInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentInfoLogic) GetCommentInfo(in *video.GetCommentInfoRequest) (*video.GetCommentInfoResponse, error) {
	var comment model.Comment
	err := l.svcCtx.Mysql.Where("comment_id = ?", in.CommentId).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrCommentNotFound.Err()
		}
		return nil, utils.InternalWithDetails("err querying comment", err)
	}

	return &video.GetCommentInfoResponse{
		Id:          comment.CommentId,
		UserId:      comment.UserId,
		Content:     comment.Content,
		CreatedTime: comment.CreatedAt.Unix(),
	}, nil
}
