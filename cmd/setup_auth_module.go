package main

import (
	grpcauth "shop/internal/api/grpc/auth"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/auth"
	service "shop/internal/service/auth"
	validator "shop/internal/validator/auth"
)

func SetupAuthModule(mysqlRepository mysql.Connection, cfg service.Config) (grpcauth.Server, grpcauth.Interceptor) {
	repo := repository.New(mysqlRepository)
	svc := service.New(repo, cfg)
	val := validator.New()
	// gin http freamwork
	// grpc server
	grpcServer := grpcauth.New(svc, val)
	grpcInterceptor := grpcauth.NewInterceptor(repo, cfg.AccessTokenSecret)

	return grpcServer, grpcInterceptor
}
