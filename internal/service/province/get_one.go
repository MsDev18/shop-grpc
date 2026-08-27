package province

import (
	"context"
	dto "shop/internal/dto/province"
)

func (s Service) GetOne(ctx context.Context, id uint) (dto.GetOneResponse, error) {
	const op = "province-service.GetOne"
	// call repository
	p, err := s.repository.GetOneByID(ctx, id)
	if err != nil {
		return dto.GetOneResponse{}, err
	}
	// map to dto.GetOneResponse
	province := dto.GetOneResponse{
		ID:   p.ID,
		Name: p.Name,
	}
	
	return province, nil
}
