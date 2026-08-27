package user

import (
	"context"
	userdto "shop/internal/dto/user"
)

func (s Service) Profile(ctx context.Context, userID uint) (userdto.ProfileResponse, error) {
	const op = "user-service.Profile"

	// get user from database
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return userdto.ProfileResponse{}, err
	}

	return userdto.ProfileResponse{
		Name:        user.Name,
		Avatar:      user.Avatar,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
