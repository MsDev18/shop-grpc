package user

import (
	"context"
	userdto "shop/internal/dto/user"
)

func (s Service) UpdateProfile(ctx context.Context, userID uint, req userdto.UpdateProfileRequest) error {
	const op = "user-service-UpdateProfile"
	// 1. process image and upload in server
	var avatarURI *string
	if req.Avatar != nil {
		uri, processImageErr := s.imageProcessor.Process(ctx , req.Avatar)
		if processImageErr != nil {
			return processImageErr
		}
		avatarURI = &uri
	}
	
	// 2. update user record in database
	// with updateProfile method
	repoErr := s.repository.UpdateProfile(ctx, userID, req.Name, avatarURI)
	if repoErr != nil {
		return repoErr
	}
	return nil
}
