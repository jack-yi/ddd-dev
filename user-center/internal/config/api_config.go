package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ApiConfig struct {
	rest.RestConf
	MySQL         MySQLConfig
	JWT           JWTConfig
	Google        GoogleConfig
	UserCenterRpc zrpc.RpcClientConf
}

type MySQLConfig struct {
	DataSource string
}

type JWTConfig struct {
	Secret string
	Expire int64
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
