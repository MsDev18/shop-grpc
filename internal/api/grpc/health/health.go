package health

import (
	pb "shop/internal/pb/health"
)

type Server struct {
	pb.UnimplementedHealthServiceServer
}

func New () Server {
	return Server{}
}


