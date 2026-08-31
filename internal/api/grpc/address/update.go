package address

import (
	"context"
	"shop/internal/api/grpc/auth"
	dto "shop/internal/dto/address"
	pb "shop/internal/pb/address"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	const op = "address-grpc"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("not found user in context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	var provinceID *uint
	if req.ProvinceId != nil {
		v := uint(req.GetProvinceId())
		provinceID = &v
	}

	dtoReq := dto.UpdateRequest{
		Title:      req.Title,
		ProvinceID: provinceID,
		City:       req.City,
		Address:    req.Address,
		PostalCode: req.PostalCode,
	}

	if validationErr := s.validator.Update(ctx , dtoReq) ; validationErr != nil {
		return nil , mapper.ErrorToGrpc(validationErr)
	}

	if serviceErr := s.service.Update(ctx, userID, uint(req.Id) , dtoReq); serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &pb.UpdateResponse{}, nil
}
