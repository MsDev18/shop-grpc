package address

import (
	"context"
	"shop/internal/api/grpc/auth"
	pb "shop/internal/pb/address"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	const op = "address-grpc.GetAll"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("user id not found in context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	list, err := s.service.GetAll(ctx, userID)
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	resp := &pb.GetAllResponse{
		Addresses: make([]*pb.AddressResponse, 0, len(list)),
	}
	for _, a := range list {
		resp.Addresses = append(resp.Addresses, &pb.AddressResponse{
			Id:         uint64(a.ID),
			Title:      a.Title,
			Province:   a.Province,
			City:       a.City,
			Address:    a.Address,
			PostalCode: a.PostalCode,
		})
	}
	return resp, nil
}
