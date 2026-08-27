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
)

func setupProductModule(config imageprocessor.Config, mysqlRepo mysql.Connection,categoryService category.Service) handler.Handler {
	repository := repository.New(mysqlRepo)

	storage := localstorage.New("product")
	imageProcessor := imageprocessor.New(config, storage)

	service := service.New(repository,categoryService, imageProcessor)
	validator := validator.New()
	return handler.New(service, validator)
}
