package config

import "github.com/zeromicro/go-zero/zrpc"

type RpcConfig struct {
	zrpc.RpcServerConf
	MySQL MySQLConfig
	JWT   JWTConfig
}
