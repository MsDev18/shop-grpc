package main

import (
	"fmt"
	grpcaddress "shop/internal/api/grpc/address"
	"shop/internal/repository/mysql"
	repository "shop/internal/repository/mysql/address"
	service "shop/internal/service/address"
	validator "shop/internal/validator/address"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	provincepb "shop/internal/pb/province"
	provinceclient "shop/internal/client/province"
)

func setupAddressModule(mysqlRepo mysql.Connection, provinceServiceAddr string) grpcaddress.Server {
	repo := repository.New(mysqlRepo)

	conn, err := grpc.NewClient(provinceServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("faild to connect province service -> %w", err))
	}

	provinceGRPCClinet := provincepb.NewProvinceServiceClient(conn)
	provinceService := provinceclient.New(provinceGRPCClinet)

	svc := service.New(repo, provinceService)
	val := validator.New()

	grpcServer := grpcaddress.New(svc, val)

	return grpcServer
}
