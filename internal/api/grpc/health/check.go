package health

import (
	"context"
	pb "shop/internal/pb/health"
)

func (s Server) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Message: "Health Check ✅",
	}, nil
}
