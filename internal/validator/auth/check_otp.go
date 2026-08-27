package auth

import (
	"context"
	"errors"
	"regexp"
	authdto "shop/internal/dto/auth"
	"shop/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) CheckOtp(ctx context.Context, req authdto.CheckOtpRequest) error {
	const op = "auth-validator.CheckOtp"
	// validate data
	err := validation.ValidateStructWithContext(ctx, &req,
		validation.Field(&req.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(IR_MOBILE_REGEX))),
		validation.Field(&req.Code, validation.Required, validation.Length(5, 5)),
	)

	if err != nil {
		var validationErr validation.Errors
		if errors.As(err, &validationErr) {
			meta := make(map[string]any)
			for key, value := range validationErr {
				meta[key] = value.Error()
			}
			return richerror.New().
				SetOp(op).
				SetMsg("inputl validation error").
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
