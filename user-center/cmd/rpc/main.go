package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/config"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"github.com/yangboyi/ddd-dev/user-center/internal/repository"
	rpcServer "github.com/yangboyi/ddd-dev/user-center/internal/server/rpc"
	"github.com/yangboyi/ddd-dev/user-center/internal/seed"
	"github.com/yangboyi/ddd-dev/user-center/proto/pb"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/rpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.RpcConfig
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

	seed.SeedRoles(context.Background(), roleRepo)

	srv := zrpc.MustNewServer(c.RpcServerConf, func(s *grpc.Server) {
		pb.RegisterUserCenterServer(s, rpcServer.NewUserCenterServer(userApp))
	})
	defer srv.Stop()

	fmt.Printf("Starting user-center-rpc at %s...\n", c.ListenOn)
	srv.Start()
}
