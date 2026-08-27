package address

import (
	"context"
	dto "shop/internal/dto/address"
)

func (s Service) GetAll(ctx context.Context, userID uint) ([]dto.CreateResponse, error) {
	const op = "address-service.GetAll"
	// call reposiotry
	provinceList, err := s.provinceService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var provinceNameByID = make(map[uint]string, len(provinceList))
	for _, value := range provinceList {
		provinceNameByID[value.ID] = value.Name
	}

	// call repository
	addressList, err := s.repository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	// map to dto.CreateResponse
	response := make([]dto.CreateResponse, len(addressList))
	for i, v := range addressList {
		response[i] = dto.CreateResponse{
			ID:         v.ID,
			Title:      v.Title,
			Province:   provinceNameByID[v.ProvinceID],
			City:       v.City,
			Address:    v.Address,
			PostalCode: v.PostalCode,
		}
	}
	return response, nil
}
