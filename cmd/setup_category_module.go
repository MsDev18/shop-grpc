package main

import (
	grpccategory "shop/internal/api/grpc/category"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/pkg/imageprocessor/localstorage"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/category"
	service "shop/internal/service/category"
	validator "shop/internal/validator/category"
)

func SetupCategoryModule(mysqlRepo mysql.Connection, uploadConfig imageprocessor.Config) (grpccategory.Server, service.Service)	 {
	repository := repository.New(mysqlRepo)
	storage := localstorage.New("category")
	imageProcessor := imageprocessor.New(uploadConfig, storage)
	svc := service.New(repository, imageProcessor)
	val := validator.New()

	categoryGrpcServer := grpccategory.New(svc, val)

	return categoryGrpcServer, svc
}
