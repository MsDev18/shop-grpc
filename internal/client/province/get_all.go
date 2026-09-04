package province

import (
	"context"
	dto "shop/internal/dto/province"
	pb "shop/internal/pb/province"
)

func (c Client) GetAll(ctx context.Context) ([]dto.GetOneResponse, error) {
	res, err := c.grpcClient.GetAll(ctx, &pb.GetAllRequest{})
	if err != nil {
		return nil, err
	}

	provinces := make([]dto.GetOneResponse, 0, len(res.Provinces))
	for _, p := range res.Provinces {
		provinces = append(provinces, dto.GetOneResponse{
			ID:   uint(p.Id),
			Name: p.Name,
		})
	}

	return provinces, nil
}
