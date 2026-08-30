package main

import (
	"shop/internal/api/grpc/province"
	handler "shop/internal/api/handler/province"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/province"
	service "shop/internal/service/province"
)

func setupProvinceModule(mysqlRepo mysql.Connection) (handler.Handler, service.Service, province.Server) {
	repository := repository.New(mysqlRepo)
	svc := service.New(repository)
	provinceGrpc := province.New(svc)
	return handler.New(svc) , svc , provinceGrpc
}
