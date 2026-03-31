package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/backend/infra/config"
	"github.com/yangboyi/ddd-dev/backend/infra/vars"
	"github.com/yangboyi/ddd-dev/backend/internal"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"github.com/yangboyi/ddd-dev/backend/internal/server"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

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
