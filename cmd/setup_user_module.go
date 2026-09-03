package main

import (
	"shop/internal/api/grpc/user"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/pkg/imageprocessor/localstorage"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/user"
	service "shop/internal/service/user"
	validator "shop/internal/validator/user"
)

func SetupUserModule(mysqlRepo mysql.Connection, uploadConfig imageprocessor.Config) user.Server {
	repository := repository.New(mysqlRepo)
	storage := localstorage.New("avatar")
	imageProcessor := imageprocessor.New(uploadConfig, storage)
	svc := service.New(repository, imageProcessor)
	val := validator.New()
	userGrpc := user.New(svc, val)
	return userGrpc
}
