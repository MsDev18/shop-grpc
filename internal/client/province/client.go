package province

import (
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
