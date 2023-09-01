package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/contact/contact"
	"tikstart/rpc/contact/internal/svc"
)

type CreateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMessageLogic {
	return &CreateMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMessageLogic) CreateMessage(in *contact.CreateMessageRequest) (*contact.Empty, error) {
	db := l.svcCtx.Mysql
	//创建并增加消息记录
	message := model.Message{
		FromId:   in.FromId,
		ToUserId: in.ToId,
		Content:  in.Content,
	}
	if err := db.Create(&message).Error; err != nil {
		return nil, utils.InternalWithDetails("error creating message", err)
	}
	return &contact.Empty{}, nil
}
