package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yangboyi/ddd-dev/backend/infra/config"
	"github.com/yangboyi/ddd-dev/backend/infra/vars"
	"github.com/yangboyi/ddd-dev/backend/internal"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"github.com/yangboyi/ddd-dev/backend/internal/server"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 如果配置了 etcd，从 etcd 覆盖配置
	if len(c.Etcd.Hosts) > 0 && c.Etcd.Key != "" {
		loadConfigFromEtcd(&c)
	}

	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql error: %v", err)
	}
	vars.DB = db

	// Auto migrate tables
	if err := db.AutoMigrate(&po.SourceItem{}, &po.Product{}, &po.ProductSKU{}, &po.PublishTask{}); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}

	srv := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer srv.Stop()

	handlers := internal.InitHandlers(db)
	server.RegisterRoutes(srv, handlers.SourceItem, handlers.Product, handlers.Publish)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	srv.Start()
}

func loadConfigFromEtcd(c *config.Config) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Etcd.Hosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("warn: connect etcd failed: %v, using file config", err)
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, c.Etcd.Key)
	if err != nil {
		log.Printf("warn: get etcd key %s failed: %v, using file config", c.Etcd.Key, err)
		return
	}

	if len(resp.Kvs) == 0 {
		log.Printf("warn: etcd key %s not found, using file config", c.Etcd.Key)
		return
	}

	// etcd 中的配置覆盖文件配置
	var etcdConfig config.Config
	if err := json.Unmarshal(resp.Kvs[0].Value, &etcdConfig); err != nil {
		log.Printf("warn: parse etcd config failed: %v, using file config", err)
		return
	}

	// 覆盖关键配置
	if etcdConfig.MySQL.DataSource != "" {
		c.MySQL.DataSource = etcdConfig.MySQL.DataSource
	}
	if etcdConfig.Host != "" {
		c.Host = etcdConfig.Host
	}
	if etcdConfig.Port > 0 {
		c.Port = etcdConfig.Port
	}

	log.Printf("config loaded from etcd: %s", c.Etcd.Key)
}
