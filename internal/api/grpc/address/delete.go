package address

import (
	"context"
	"shop/internal/api/grpc/auth"
	pb "shop/internal/pb/address"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	const op = "address-grpc.Delete"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("not found user id in context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	if serviceErr := s.service.Delete(ctx, userID, uint(req.Id)); serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &pb.DeleteResponse{}, nil
}
