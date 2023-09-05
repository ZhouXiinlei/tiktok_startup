package svc

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"tikstart/common/tikcos"
	"tikstart/http/internal/config"
	"tikstart/http/internal/middleware"
	"tikstart/rpc/contact/contactClient"
	"tikstart/rpc/user/userClient"
	"tikstart/rpc/video/videoClient"
)

type ServiceContext struct {
	Config           config.Config
	UserRpc          userClient.User
	VideoRpc         videoClient.Video
	ContactRpc       contactClient.Contact
	JwtAuth          rest.Middleware
	TengxunyunClient *cos.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		UserRpc:          userClient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:         videoClient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		ContactRpc:       contactClient.NewContact(zrpc.MustNewClient(c.ContactRpc)),
		TengxunyunClient: tikcos.TengxunyunInit(c.COS),
		JwtAuth:          middleware.NewJwtAuthMiddleware(c).Handle,
	}
}
