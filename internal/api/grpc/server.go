package grpc

import (
	"fmt"
	"log"
	"net"
	"shop/internal/api/grpc/auth"
	"shop/internal/api/grpc/health"
	healthpb "shop/internal/pb/health"
	authpb "shop/internal/pb/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer   *grpc.Server
	healthServer health.Server
	authServer auth.Server
	address      string
}


func New (address string, healthServer health.Server, authServer auth.Server, authInterceptor auth.Interceptor) Server {
	return Server{
		address:      address,
		grpcServer:   grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary())),
		healthServer: healthServer,
		authServer: authServer,
	}
}

func (s Server) Run () {
	healthpb.RegisterHealthServiceServer(s.grpcServer , s.healthServer)
	authpb.RegisterAuthServiceServer(s.grpcServer , s.authServer)

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