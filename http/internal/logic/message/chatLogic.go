package message

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"strconv"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/rpc/contact/contact"
)

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(req *types.MessageChatRequest) (resp *types.MessageChatResponse, err error) {
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	var preTime int64
	if req.PreMsgTime > 1e12 {
		preTime = req.PreMsgTime / 1000
	} else {
		preTime = req.PreMsgTime
	}

	res, err := l.svcCtx.ContactRpc.GetMessageList(l.ctx, &contact.GetMessageListRequest{
		FromId:     req.ToUserId,
		ToId:       userClaims.UserId,
		PreMsgTime: preTime,
	})
	if err != nil {
		return nil, utils.ReturnInternalError(l.ctx, status.Convert(err), err)
	}

	messageList := make([]types.Message, 0, len(res.Messages))
	for _, message := range res.Messages {
		messageList = append(messageList, types.Message{
			Id:         message.MessageId,
			ToUserId:   message.ToId,
			FromUserId: message.FromId,
			Content:    message.Content,
			CreateTime: strconv.FormatInt(message.CreateTime, 10),
		})
	}

	return &types.MessageChatResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		MessageList: messageList,
	}, nil
}
