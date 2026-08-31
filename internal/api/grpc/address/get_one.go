package address

import (
	"context"
	"shop/internal/api/grpc/auth"
	pb "shop/internal/pb/address"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) GetOne(ctx context.Context, req *pb.GetOneRequest) (*pb.AddressResponse, error) {
	const op = "address-grpc.GetOne"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("not found user id from context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	res, serviceErr := s.service.GetOne(ctx, userID, uint(req.GetId()))
	if serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &pb.AddressResponse{
		Id:         uint64(res.ID),
		Title:      res.Title,
		Province:   res.Province,
		City:       res.City,
		Address:    res.Address,
		PostalCode: res.PostalCode,
	}, nil
}
