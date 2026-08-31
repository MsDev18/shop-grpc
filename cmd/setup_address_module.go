package main

import (
	grpcaddress "shop/internal/api/grpc/address"
	handler "shop/internal/api/handler/address"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/address"
	service "shop/internal/service/address"
	"shop/internal/service/province"
	validator "shop/internal/validator/address"
)

func setupAddressModule(mysqlRepo mysql.Connection, provinceService province.Service) (handler.Handler, grpcaddress.Server) {
	repo := repository.New(mysqlRepo)
	svc := service.New(repo, provinceService)
	val := validator.New()

	grpcServer := grpcaddress.New(svc, val)

	return handler.New(svc, val), grpcServer
}
