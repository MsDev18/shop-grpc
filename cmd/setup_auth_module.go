package main

import (
	handler "shop/internal/api/handler/auth"
	middleware "shop/internal/api/middleware/auth"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/auth"
	service "shop/internal/service/auth"
	validator "shop/internal/validator/auth"
	grpcauth "shop/internal/api/grpc/auth"
)

func SetupAuthModule(mysqlRepository mysql.Connection, cfg service.Config) (handler.Handler, middleware.Middleware, grpcauth.Server, grpcauth.Interceptor) {
	repo := repository.New(mysqlRepository)
	svc := service.New(repo, cfg)
	val := validator.New()
	// gin http freamwork 
	handler := handler.New(svc, val)
	middleware := middleware.New(repo , cfg.AccessTokenSecret)
	// grpc server
	grpcServer := grpcauth.New(svc , val)
	grpcInterceptor := grpcauth.NewInterceptor(repo , cfg.AccessTokenSecret)

	return handler , middleware , grpcServer , grpcInterceptor
}