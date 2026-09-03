package main

import (
	handler "shop/internal/api/handler/product"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/pkg/imageprocessor/localstorage"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/product"
	"shop/internal/service/category"
	service "shop/internal/service/product"
	validator "shop/internal/validator/product"
	grpcproduct "shop/internal/api/grpc/product"
)

func setupProductModule(config imageprocessor.Config, mysqlRepo mysql.Connection,categoryService category.Service) (handler.Handler, grpcproduct.Server) {
	repository := repository.New(mysqlRepo)

	storage := localstorage.New("product")
	imageProcessor := imageprocessor.New(config, storage)

	svc := service.New(repository,categoryService, imageProcessor)
	val := validator.New()

	grpcproduct := grpcproduct.New(svc , val)

	return handler.New(svc, val), grpcproduct
}
