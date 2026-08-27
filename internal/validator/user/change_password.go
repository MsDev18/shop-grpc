package user

import (
	"context"
	"errors"
	userdto "shop/internal/dto/user"
	"shop/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ChangePassword(ctx context.Context, req userdto.ChangePasswordRequest) error {
	const op = "user-validator.ChangePassword"
	err := validation.ValidateStructWithContext(
		ctx,
		&req,
		validation.Field(&req.Password, validation.Required, validation.Length(8, 0)),
		validation.Field(&req.ConfirmPassword, validation.Required, validation.Length(8, 0)),
	)
	if err != nil {
		var validationErr validation.Errors
		if errors.As(err, &validationErr) {
			meta := make(map[string]any)
			for field, err := range validationErr {
				meta[field] = err.Error()
			}
			return richerror.New().
				SetOp(op).
				SetMsg("input invalid").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err).
				SetMeta(meta)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return nil
}
