package grpc

import (
	"fmt"
	"log"
	"net"
	"shop/internal/api/grpc/auth"
	"shop/internal/api/grpc/category"
	"shop/internal/api/grpc/health"
	"shop/internal/api/grpc/province"
	"shop/internal/api/grpc/user"
	authpb "shop/internal/pb/auth"
	categorypb "shop/internal/pb/category"
	healthpb "shop/internal/pb/health"
	provincepb "shop/internal/pb/province"
	userpb "shop/internal/pb/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer     *grpc.Server
	healthServer   health.Server
	authServer     auth.Server
	userServer     user.Server
	categoryServer category.Server
	provinceServer province.Server
	address        string
}

func New(address string, healthServer health.Server, authServer auth.Server, userServer user.Server, categoryServer category.Server, provinceServer province.Server, authInterceptor auth.Interceptor) Server {
	return Server{
		address:        address,
		grpcServer:     grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary())),
		healthServer:   healthServer,
		authServer:     authServer,
		userServer:     userServer,
		categoryServer: categoryServer,
		provinceServer: provinceServer,
	}
}

func (s Server) Run() {
	healthpb.RegisterHealthServiceServer(s.grpcServer, s.healthServer)
	authpb.RegisterAuthServiceServer(s.grpcServer, s.authServer)
	userpb.RegisterUserServiceServer(s.grpcServer, s.userServer)
	categorypb.RegisterCategoryServiceServer(s.grpcServer, s.categoryServer)
	provincepb.RegisterProvinceServiceServer(s.grpcServer, s.provinceServer)

	reflection.Register(s.grpcServer)

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(fmt.Errorf("error on listen server on port %s", s.address))
	}

	log.Printf("GRPC server listening on %s", s.address)
	if err := s.grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
