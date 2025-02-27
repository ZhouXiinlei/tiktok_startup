package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common"
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
	var count int64
	err := l.svcCtx.DB.Model(&model.User{}).Where("user_id = ?", in.ToId).Count(&count).Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying user", err)
	}
	if count == 0 {
		return nil, common.ErrUserNotFound.Err()
	}

	message := model.Message{
		FromId:   in.FromId,
		ToUserId: in.ToId,
		Content:  in.Content,
	}
	if err := l.svcCtx.DB.Create(&message).Error; err != nil {
		return nil, utils.InternalWithDetails("error creating message", err)
	}
	return &contact.Empty{}, nil
}
