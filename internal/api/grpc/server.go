package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"shop/internal/api/grpc/address"
	"shop/internal/api/grpc/auth"
	"shop/internal/api/grpc/category"
	"shop/internal/api/grpc/product"
	"shop/internal/api/grpc/user"
	addresspb "shop/internal/pb/address"
	authpb "shop/internal/pb/auth"
	categorypb "shop/internal/pb/category"
	productpb "shop/internal/pb/product"
	userpb "shop/internal/pb/user"
)

type Server struct {
	grpcServer     *grpc.Server
	authServer     auth.Server
	userServer     user.Server
	categoryServer category.Server
	addressServer  address.Server
	productServer  product.Server
	address        string
}

func New(addr string, authServer auth.Server, userServer user.Server, categoryServer category.Server, addressServer address.Server, productServer product.Server, authInterceptor auth.Interceptor) Server {
	roleInterceptor := auth.NewRoleInterceptor()
	return Server{
		address: addr,
		grpcServer: grpc.NewServer(
			grpc.ChainUnaryInterceptor(authInterceptor.Unary(), roleInterceptor.Unary()),
			grpc.ChainStreamInterceptor(authInterceptor.Stream(), roleInterceptor.Stream()),
		),
		authServer:     authServer,
		userServer:     userServer,
		categoryServer: categoryServer,
		addressServer:  addressServer,
		productServer:  productServer,
	}
}

func (s Server) Run() {
	authpb.RegisterAuthServiceServer(s.grpcServer, s.authServer)
	userpb.RegisterUserServiceServer(s.grpcServer, s.userServer)
	categorypb.RegisterCategoryServiceServer(s.grpcServer, s.categoryServer)
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
