package main

import (
	handler "shop/internal/api/handler/category"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/pkg/imageprocessor/localstorage"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/category"
	service "shop/internal/service/category"
	validator "shop/internal/validator/category"
	grpccategory "shop/internal/api/grpc/category"
)
	

func SetupCategoryModule (mysqlRepo mysql.Connection, uploadConfig imageprocessor.Config) (handler.Handler , service.Service, grpccategory.Server){
	repository := repository.New(mysqlRepo)
	storage := localstorage.New("category")
	imageProcessor := imageprocessor.New(uploadConfig , storage)
	svc := service.New(repository , imageProcessor)
	val := validator.New() 

	categoryGrpcServer := grpccategory.New(svc , val)

	handler := handler.New(svc, val)
	return handler, svc , categoryGrpcServer
}