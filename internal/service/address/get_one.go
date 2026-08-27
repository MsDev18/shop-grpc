package address

import (
	"context"
	dto "shop/internal/dto/address"
)

func (s Service) GetOne(ctx context.Context, userID uint, addressID uint) (dto.CreateResponse, error) {
	const op = "address-service.GetOne"

	// get address form repository
	a, err := s.repository.GetOne(ctx, userID, addressID)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	// convert address.provinceID to province name
	province, err := s.provinceService.GetOne(ctx, a.ProvinceID)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	// map to dto.CreateResponse
	res := dto.CreateResponse{
		ID:         a.ID,
		Title:      a.Title,
		Province:   province.Name,
		City:       a.City,
		Address:    a.Address,
		PostalCode: a.PostalCode,
	}
	// return
	return res, nil
}
