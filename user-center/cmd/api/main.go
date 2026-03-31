package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/config"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"github.com/yangboyi/ddd-dev/user-center/internal/repository"
	apiServer "github.com/yangboyi/ddd-dev/user-center/internal/server/api"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.ApiConfig
	conf.MustLoad(*configFile, &c)

	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql error: %v", err)
	}

	if err := db.AutoMigrate(&po.User{}, &po.Role{}, &po.UserRole{}); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}

	userRepo := repository.NewUserRepoImpl(db)
	roleRepo := repository.NewRoleRepoImpl(db)
	jwtCfg := application.JWTConfig{Secret: c.JWT.Secret, Expire: c.JWT.Expire}
	userApp := application.NewUserApp(userRepo, roleRepo, jwtCfg, db)

	authHandler := apiServer.NewAuthHandler(userApp, c.Google)
	userHandler := apiServer.NewUserHandler(userApp)
	initHandler := apiServer.NewInitHandler(userApp)

	srv := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer srv.Stop()

	apiServer.RegisterRoutes(srv, authHandler, userHandler, initHandler, c.JWT.Secret)

	fmt.Printf("Starting user-center-api at %s:%d...\n", c.Host, c.Port)
	srv.Start()
}
