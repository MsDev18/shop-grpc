package province

import (
	"context"
	dto "shop/internal/dto/province"
)

func (s Service) GetAll(ctx context.Context) ([]dto.GetOneResponse, error) {
	const op = "province-service.GetAll"

	// call repository
	provinces, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	// map to dto.GetOneResponse
	var response = make([]dto.GetOneResponse, len(provinces))
	for index, value := range provinces {
		response[index] = dto.GetOneResponse{
			ID:   value.ID,
			Name: value.Name,
		}
	}
	// return response 
	return response, nil
}
