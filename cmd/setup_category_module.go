package main

import (
	handler "shop/internal/api/handler/category"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/pkg/imageprocessor/localstorage"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/category"
	service "shop/internal/service/category"
	validator "shop/internal/validator/category"
)
	

func SetupCategoryModule (mysqlRepo mysql.Connection, uploadConfig imageprocessor.Config) (handler.Handler , service.Service){
	repository := repository.New(mysqlRepo)
	storage := localstorage.New("category")
	imageProcessor := imageprocessor.New(uploadConfig , storage)
	service := service.New(repository , imageProcessor)
	validator := validator.New() 
	handler := handler.New(service, validator)
	return handler, service
}