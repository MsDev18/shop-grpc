package main

import (
	_ "github.com/go-sql-driver/mysql"
	grpcserver "shop/internal/api/grpc"
	grpchealth "shop/internal/api/grpc/health"
	"shop/internal/api/handler/health"
	"shop/internal/api/server"
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
	healthHandler := health.New()
	authHandler, authMiddleware, grpcAuthServer, authInterceptor := SetupAuthModule(mysqlRepo, config.AuthService)
	userHandler, grpcUserServer := SetupUserModule(mysqlRepo, config.Upload)
	categoryHandler, categoryService, grpcCategoryServer := SetupCategoryModule(mysqlRepo, config.Upload)
	provinceHandler, provinceService := setupProvinceModule(mysqlRepo)
	addressHandler := setupAddressModule(mysqlRepo, provinceService)
	productHandler := setupProductModule(config.Upload, mysqlRepo, categoryService)

	grpcHealthServer := grpchealth.New()

	go grpcserver.New(
		"0.0.0.0:50051",
		grpcHealthServer,
		grpcAuthServer,
		grpcUserServer,
		grpcCategoryServer,
		authInterceptor,
	).Run()

	// create new http server and run it
	httpServer := server.New(
		config.Server,
		healthHandler,
		authHandler,
		userHandler,
		categoryHandler,
		provinceHandler,
		addressHandler,
		productHandler,
		authMiddleware,
	)
	httpServer.Run()
}
