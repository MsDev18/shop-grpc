package province

import (
	"context"
	dto "shop/internal/dto/province"
	pb "shop/internal/pb/province"
)

type Client struct {
	grpcClient pb.ProvinceServiceClient
}

func New(grpcClient pb.ProvinceServiceClient) Client {
	return Client{
		grpcClient: grpcClient,
	}
}

func (c Client) GetOne(ctx context.Context, id uint) (dto.GetOneResponse, error) {
	res, err := c.grpcClient.GetOne(ctx, &pb.GetOneRequest{Id: uint64(id)})
	if err != nil {
		return dto.GetOneResponse{}, err
	}
	return dto.GetOneResponse{
		ID:   uint(res.Id),
		Name: res.Name,
	}, nil
}
