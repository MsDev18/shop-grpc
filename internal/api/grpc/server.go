package grpc

import (
	"fmt"
	"log"
	"net"
	"shop/internal/api/grpc/address"
	"shop/internal/api/grpc/auth"
	"shop/internal/api/grpc/category"
	"shop/internal/api/grpc/health"
	"shop/internal/api/grpc/product"
	"shop/internal/api/grpc/province"
	"shop/internal/api/grpc/user"
	addresspb "shop/internal/pb/address"
	authpb "shop/internal/pb/auth"
	categorypb "shop/internal/pb/category"
	healthpb "shop/internal/pb/health"
	provincepb "shop/internal/pb/province"
	userpb "shop/internal/pb/user"
	productpb "shop/internal/pb/product"
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
	addressServer  address.Server
	productServer product.Server
	address        string
}

func New(addr string, healthServer health.Server, authServer auth.Server, userServer user.Server, categoryServer category.Server, provinceServer province.Server, addressServer address.Server,productServer product.Server, authInterceptor auth.Interceptor) Server {
	roleInterceptor := auth.NewRoleInterceptor()
	return Server{
		address:        addr,
		grpcServer:     grpc.NewServer(
			grpc.ChainUnaryInterceptor(authInterceptor.Unary(), roleInterceptor.Unary()),
			grpc.ChainStreamInterceptor(authInterceptor.Stream(), roleInterceptor.Stream()),
		),
		healthServer:   healthServer,
		authServer:     authServer,
		userServer:     userServer,
		categoryServer: categoryServer,
		provinceServer: provinceServer,
		addressServer:  addressServer,
		productServer:  productServer,	
	}
}

func (s Server) Run() {
	healthpb.RegisterHealthServiceServer(s.grpcServer, s.healthServer)
	authpb.RegisterAuthServiceServer(s.grpcServer, s.authServer)
	userpb.RegisterUserServiceServer(s.grpcServer, s.userServer)
	categorypb.RegisterCategoryServiceServer(s.grpcServer, s.categoryServer)
	provincepb.RegisterProvinceServiceServer(s.grpcServer, s.provinceServer)
	addresspb.RegisterAddressServiceServer(s.grpcServer, s.addressServer)
	productpb.RegisterProductServiceServer(s.grpcServer, s.productServer)

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
