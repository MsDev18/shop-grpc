package main

import (
	handler "shop/internal/api/handler/auth"
	middleware "shop/internal/api/middleware/auth"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/auth"
	service "shop/internal/service/auth"
	validator "shop/internal/validator/auth"
)

func SetupAuthModule(mysqlRepository mysql.Connection, cfg service.Config) (handler.Handler, middleware.Middleware) {
	repository := repository.New(mysqlRepository)
	service := service.New(repository, cfg)
	validator := validator.New()
	handler := handler.New(service, validator)
	middleware := middleware.New(repository , cfg.AccessTokenSecret)
	return handler , middleware
}