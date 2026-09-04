package main

import (
	_ "github.com/go-sql-driver/mysql"
	grpcserver "shop/internal/api/grpc"
	"shop/internal/migrator"
	"shop/internal/repository/mysql"
)

func main() {
	// load project configuration
	config := LoadConfig()
	// migration
	m := migrator.New(config.MySQL.GetDSN())
	if mErr := m.Up(); mErr != nil {
		panic(mErr)
	}
	// mysql pure connection
	mysqlRepo := mysql.New(config.MySQL)

	// setup project handlers
	grpcAuthServer, authInterceptor := SetupAuthModule(mysqlRepo, config.AuthService)
	grpcUserServer := SetupUserModule(mysqlRepo, config.Upload)
	grpcCategoryServer, categoryService := SetupCategoryModule(mysqlRepo, config.Upload)
	grpcAddressServer := setupAddressModule(mysqlRepo, "0.0.0.0:50053")
	grpcProductServer := setupProductModule(config.Upload, mysqlRepo, categoryService)


	grpcserver.New(
		"0.0.0.0:50051",
		grpcAuthServer,
		grpcUserServer,
		grpcCategoryServer,
		grpcAddressServer,
		grpcProductServer,
		authInterceptor,
	).Run()
}
