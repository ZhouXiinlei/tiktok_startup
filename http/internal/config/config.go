package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"tikstart/common/tikcos"
)

type Config struct {
	rest.RestConf
	UserRpc    zrpc.RpcClientConf
	VideoRpc   zrpc.RpcClientConf
	ContactRpc zrpc.RpcClientConf
	JwtAuth    struct {
		Secret string
		Expire int64
	}
	COS        tikcos.TengxunyunCfg
	CDNBaseURL string
}
