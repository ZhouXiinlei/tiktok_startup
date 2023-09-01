package video

import (
	"context"

	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return
}
