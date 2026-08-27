package main

import (
	handler "shop/internal/api/handler/province"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/province"
	service "shop/internal/service/province"
)

func setupProvinceModule(mysqlRepo mysql.Connection) (handler.Handler, service.Service) {
	repository := repository.New(mysqlRepo)
	service := service.New(repository)
	return handler.New(service) , service
}
