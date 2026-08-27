package main

import (
	handler "shop/internal/api/handler/user"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/pkg/imageprocessor/localstorage"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/user"
	service "shop/internal/service/user"
	validator "shop/internal/validator/user"
)

func SetupUserModule(mysqlRepo  mysql.Connection, uploadConfig imageprocessor.Config) handler.Handler {
	repository := repository.New(mysqlRepo)
	storage := localstorage.New("avatar")
	imageProcessor := imageprocessor.New(uploadConfig , storage)
	service  := service.New(repository, imageProcessor)
	validator := validator.New()
	handler := handler.New(service , validator)
	return handler
}