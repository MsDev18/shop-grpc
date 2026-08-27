package main

import (
	"shop/internal/api/handler/health"
	"shop/internal/api/server"
	"shop/internal/migrator"
	"shop/internal/repository/mysql"

	_ "github.com/go-sql-driver/mysql"
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
	authHandler, authMiddleware := SetupAuthModule(mysqlRepo, config.AuthService)
	userHandler := SetupUserModule(mysqlRepo, config.Upload)
	categoryHandler, categoryService := SetupCategoryModule(mysqlRepo, config.Upload)
	provinceHandler, provinceService := setupProvinceModule(mysqlRepo)
	addressHandler := setupAddressModule(mysqlRepo, provinceService)
	productHandler := setupProductModule(config.Upload , mysqlRepo, categoryService)
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
