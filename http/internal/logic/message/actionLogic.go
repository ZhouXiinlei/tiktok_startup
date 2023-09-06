package message

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strings"
	"tikstart/common"
	"tikstart/common/utils"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	"tikstart/rpc/contact/contact"
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

	regPattern := regexp.MustCompile("\\s+")
	content := regPattern.ReplaceAllString(req.Content, "")
	content = strings.Trim(content, " ")
	if content == "" {
		return nil, schema.ApiError{
			StatusCode: 422,
			Code:       42206,
			Message:    "消息内容不能为空",
		}
	}

	// no need to handle error because we have middleware to intercept it
	userClaims, _ := utils.ParseToken(req.Token, l.svcCtx.Config.JwtAuth.Secret)

	_, err = l.svcCtx.ContactRpc.CreateMessage(l.ctx, &contact.CreateMessageRequest{
		FromId:  userClaims.UserId,
		ToId:    req.ToUserId,
		Content: content,
	})
	if err != nil {
		if st, match := utils.MatchError(err, common.ErrUserNotFound); match {
			return nil, schema.ApiError{
				StatusCode: 422,
				Code:       42202,
				Message:    "用户不存在",
			}
		} else {
			return nil, utils.ReturnInternalError(l.ctx, st, err)
		}
	}

	return &types.MessageActionResponse{}, nil
}
