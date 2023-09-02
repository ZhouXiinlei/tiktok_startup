package message

import (
	"context"
	"google.golang.org/grpc/status"
	"strconv"
	"tikstart/common/utils"
	"tikstart/rpc/contact/contact"

	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

	res, err := l.svcCtx.ContactRpc.GetMessageList(l.ctx, &contact.GetMessageListRequest{
		FromId:     userClaims.UserId,
		ToId:       req.ToUserId,
		PreMsgTime: req.PreMsgTime,
	})
	if err != nil {
		st, _ := status.FromError(err)
		return nil, utils.ReturnInternalError(st, err)
	}

	messageList := make([]types.Message, 0, len(res.Messages))
	for _, message := range res.Messages {
		messageList = append(messageList, types.Message{
			Id:         message.Id,
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
