package main

import (
	"shop/internal/api/grpc/province"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/province"
	service "shop/internal/service/province"
)

func setupProvinceModule(mysqlRepo mysql.Connection) (province.Server, service.Service) {
	repository := repository.New(mysqlRepo)
	svc := service.New(repository)
	provinceGrpc := province.New(svc)
	return provinceGrpc, svc
}
