package user

import (
	"context"
	userdto "shop/internal/dto/user"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)

func (s Service) ChangePassword(ctx context.Context, userID uint, req userdto.ChangePasswordRequest) error {
	const op = "user-service.ChangePassword"

	if req.Password != req.ConfirmPassword {
		return richerror.New().
			SetOp(op).
			SetMsg("password and confirm password do not match").
			SetKind(richerror.KindBadRequestErr)
	}
	// hash password
	hashedPassBytes, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("failed to hash password").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(hashErr)
	}
	hashPass := string(hashedPassBytes)
	// create user model
	user := entity.User{
		ID:       userID,
		Password: hashPass,
	}
	// save in repository
	err := s.repository.UpdatePassword(ctx , user)
	if err != nil {
		return err
	}
	// return response
	return nil
}
