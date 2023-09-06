package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/common/utils"
	"time"

	"tikstart/rpc/contact/contact"
	"tikstart/rpc/contact/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageListLogic) GetMessageList(in *contact.GetMessageListRequest) (*contact.GetMessageListResponse, error) {
	var messages []model.Message

	if in.PreMsgTime == 0 {
		// getting history messages
		err := l.svcCtx.DB.
			Where("from_id = ? AND to_user_id = ?", in.FromId, in.ToId).
			Or("from_id = ? AND to_user_id = ?", in.ToId, in.FromId).
			Order("created_at ASC").
			Find(&messages).Error
		if err != nil {
			return nil, utils.InternalWithDetails("error querying history messages", err)
		}
	} else {
		// getting latest messages, +1 to offset millisecond
		err := l.svcCtx.DB.Where("from_id = ? AND to_user_id = ? AND created_at > ?", in.FromId, in.ToId, time.Unix(in.PreMsgTime+1, 0)).Find(&messages).Error
		if err != nil {
			return nil, utils.InternalWithDetails("error querying latest messages", err)
		}
	}

	var messageList []*contact.Message
	for _, message := range messages {
		messageList = append(messageList, &contact.Message{
			MessageId:  message.MessageId,
			Content:    message.Content,
			CreateTime: message.CreatedAt.Unix(),
			FromId:     message.FromId,
			ToId:       message.ToUserId,
		})
	}
	return &contact.GetMessageListResponse{
		Messages: messageList,
	}, nil
}
