package address

import (
	"context"
	"shop/internal/api/grpc/auth"
	pb "shop/internal/pb/address"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
	dto "shop/internal/dto/address"
)

func (s Server) Create (ctx context.Context, req *pb.CreateRequest) (*pb.AddressResponse, error) {
	const op = "address-grpc.Create"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("user id not found in request context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	dtoReq := dto.CreateRequest{
		Title:      req.GetTitle(),
		ProvinceID: uint(req.GetProvinceId()),
		City:       req.GetCity(),
		Address:    req.GetAddress(),
		PostalCode: req.GetPostalCode(),
	}

	if validationErr := s.validator.Create(ctx , dtoReq) ; validationErr != nil {
		return nil , mapper.ErrorToGrpc(validationErr)
	}

	res, serviceErr := s.service.Create(ctx, userID, dtoReq)
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