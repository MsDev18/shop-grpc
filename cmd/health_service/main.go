package main

import (
	"fmt"
	"log"
	"net"
	"shop/internal/api/grpc/health"
	healthpb "shop/internal/pb/health"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	const addr = "0.0.0.0:50052"

	healthServer := health.New()

	grpcServer := grpc.NewServer()
	healthpb.RegisterHealthServiceServer(grpcServer, healthServer)
	reflection.Register(grpcServer)

	listener ,err := net.Listen("tcp", addr)
	if err != nil{
		panic(fmt.Errorf("error on listen server on port %s", addr))
	}

	log.Printf("health service listening on %s", addr)
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
