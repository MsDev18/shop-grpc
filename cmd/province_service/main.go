package main

import (
	"fmt"
	"log"
	"net"

	grpcprovince "shop/internal/api/grpc/province"
	"shop/internal/config"
	"shop/internal/migrator"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/province"
	"shop/internal/service/province"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	provincepb "shop/internal/pb/province"
)

func main() {
	const addr = "0.0.0.0:50053"

	// this service's own config, loaded from its own .env — completely independent
	// from the main monolith's .env
	appConfig := config.New()
	appConfig.LoadFromDotEnv("cmd/province_service/.env")
	mysqlConfig := appConfig.GetMySQLConfig()

	// this service's own database migration, against its own database
	m := migrator.NewWithSourcePath("migrations/province", mysqlConfig.GetDSN())
	if err := m.Up(); err != nil {
		panic(err)
	}

	mysqlRepo := mysql.New(mysqlConfig)

	repo := repository.New(mysqlRepo)
	svc := province.New(repo)
	provinceServer := grpcprovince.New(svc)

	grpcServer := grpc.NewServer()
	provincepb.RegisterProvinceServiceServer(grpcServer, provinceServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("error on listen server on port %s", addr))
	}

	log.Printf("province service listening on %s", addr)
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}