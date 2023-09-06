package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	MySQL struct {
		Host        string
		Port        int
		User        string
		Password    string
		Database    string
		TablePrefix string
	}
	Redis struct {
		Addr string
		Pass string
		DB   int
	}
	VideoRpc   zrpc.RpcClientConf
	UserRpc    zrpc.RpcClientConf
	ContactRpc zrpc.RpcClientConf
}
