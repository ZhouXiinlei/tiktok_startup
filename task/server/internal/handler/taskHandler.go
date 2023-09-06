package handler

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/task/server/internal/svc"
)

type TaskHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskHandler(ctx context.Context, svcCtx *svc.ServiceContext) *TaskHandler {
	return &TaskHandler{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
