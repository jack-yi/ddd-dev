package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySQL MySQLConfig
	Etcd  EtcdConfig `json:",optional"`
}

type MySQLConfig struct {
	DataSource string
}

type EtcdConfig struct {
	Hosts []string `json:",optional"`
	Key   string   `json:",optional"`
}
