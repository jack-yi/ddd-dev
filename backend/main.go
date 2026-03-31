package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/backend/infra/config"
	"github.com/yangboyi/ddd-dev/backend/infra/vars"
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

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
