package video

import (
	"context"
	"fmt"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVideoLogic {
	return &PublishVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishVideoLogic) PublishVideo(req *types.PublishVideoRequest) (resp *types.PublishVideoResponse, err error) {

	//logx.WithContext(l.ctx).Infof("发布视频: %v", req)
	//resp.BasicResponse = types.BasicResponse{
	//	StatusCode: 0,
	//	StatusMsg:  "Success",
	//}
	fmt.Println(req.Title)

	return &types.PublishVideoResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
	}, nil
}
