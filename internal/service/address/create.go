package address

import (
	"context"
	dto "shop/internal/dto/address"
	"shop/internal/entity"
)

func (s Service) Create(ctx context.Context, userID uint, req dto.CreateRequest) (dto.CreateResponse, error) {
	const op = "address-service.Create"
	// check exists province id
	province, err := s.provinceService.GetOne(ctx, req.ProvinceID)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	// map dto.CreateRequest to entity.Address
	address := entity.Address{
		UserID:     userID,
		Title:      req.Title,
		ProvinceID: province.ID,
		City:       req.City,
		Address:    req.Address,
		PostalCode: req.PostalCode,
	}
	// call repository method
	a, err := s.repository.Create(ctx, address)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	// map repository response to dro.CreateResponse
	response := dto.CreateResponse{
		ID:         a.ID,
		Title:      a.Title,
		Province:   province.Name,
		City:       a.City,
		Address:    a.Address,
		PostalCode: a.PostalCode,
	}
	// return respnse
	return response , nil
}
