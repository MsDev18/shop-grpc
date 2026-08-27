package address

import (
	"context"
	dto "shop/internal/dto/address"
)

func (s Service) Update(ctx context.Context, userID uint, addressID uint, req dto.UpdateRequest) error {
	const op = "address-service.Update"
	// check exists category
	a, err := s.repository.GetOne(ctx, userID, addressID)
	if err != nil {
		return err
	}
	// check exists province
	if req.ProvinceID != nil {
		_, err := s.provinceService.GetOne(ctx, *req.ProvinceID)
		if err != nil {
			return err
		}
	}
	// call repository
	err = s.repository.Update(ctx, a.UserID, a.ID, req.Title, req.ProvinceID, req.City, req.Address, req.PostalCode)
	if err != nil {
		return err
	}
	return nil
}
