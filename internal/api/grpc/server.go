package grpc

import (
	"fmt"
	"log"
	"net"
	"shop/internal/api/grpc/health"
	healthpb "shop/internal/pb/health"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer   *grpc.Server
	healthServer health.Server
	address      string
}


func New (address string, healthServer health.Server) Server {
	return Server{
		grpcServer:   grpc.NewServer(),
		healthServer: healthServer,
		address:      address,
	}
}

func (s Server) Run () {
	healthpb.RegisterHealthServiceServer(s.grpcServer , s.healthServer)

	reflection.Register(s.grpcServer)

	listener , err := net.Listen("tcp",s.address)
	if err != nil {
		panic(fmt.Errorf("error on listen server on port %s" , s.address))
	}

	log.Printf("GRPC server listening on %s" , s.address)
	if err := s.grpcServer.Serve(listener) ; err != nil {
		panic(err)
	}
}