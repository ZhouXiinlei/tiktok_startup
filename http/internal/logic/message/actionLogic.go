package message

import (
	"context"
	"google.golang.org/grpc/status"
	"tikstart/common/utils"
	"tikstart/http/schema"
	"tikstart/rpc/contact/contact"

	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActionLogic) Action(req *types.MessageActionRequest) (resp *types.MessageActionResponse, err error) {
	if req.ActionType != 1 {
		return nil, schema.ApiError{
			StatusCode: 422,
			Code:       42205,
			Message:    "未知操作",
		}
	}

	// no need to handle error because we have middleware to intercept it
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	_, err = l.svcCtx.ContactRpc.CreateMessage(l.ctx, &contact.CreateMessageRequest{
		FromId:  userClaims.UserId,
		ToId:    req.ToUserId,
		Content: req.Content,
	})

	// If error occurred, then it's an internal error.
	if err != nil {
		st, _ := status.FromError(err)
		return nil, utils.ReturnInternalError(st, err)
	}

	return &types.MessageActionResponse{}, nil
}
