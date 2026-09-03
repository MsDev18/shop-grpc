package main

import (
	grpcaddress "shop/internal/api/grpc/address"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/address"
	service "shop/internal/service/address"
	"shop/internal/service/province"
	validator "shop/internal/validator/address"
)

func setupAddressModule(mysqlRepo mysql.Connection, provinceService province.Service) ( grpcaddress.Server) {
	repo := repository.New(mysqlRepo)
	svc := service.New(repo, provinceService)
	val := validator.New()

	grpcServer := grpcaddress.New(svc, val)

	return  grpcServer
}
