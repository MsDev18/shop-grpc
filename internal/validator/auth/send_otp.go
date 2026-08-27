package auth

import (
	"context"
	"errors"
	"regexp"
	authdto "shop/internal/dto/auth"
	"shop/internal/pkg/richerror"

	"github.com/go-ozzo/ozzo-validation/v4"
)

const (
	IR_MOBILE_REGEX = `^09[0-9]{9}$`
)

// we need to validation uniqueness phone number ?
// no - because in business logic
// if not exist user - create user with base data
// if exist user - update otp and send
func (v Validator) SendOtp(ctx context.Context, req authdto.SendOtpRequest) (error) {
	const op = "auth-validator.SendOtp"
	err := validation.ValidateStructWithContext(ctx, &req,
		// validation phone number format
		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(IR_MOBILE_REGEX)),
		),
	)

	// error handling ozzo-vlidation package
	if err != nil {
		// check type error is validation.Errors
		var validationErr validation.Errors
		// in this step validation error map[string]error
		if errors.As(err, &validationErr) {
			// create meta data for richerror package
			meta := make(map[string]any)
			for key, value := range validationErr {
				meta[key] = value.Error()
			}
			// return error
			return richerror.New().
				SetOp(op).
				SetMsg("input validation error").
				SetKind(richerror.KindBadRequestErr).
				SetMeta(meta)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected errror").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return  nil
}
// ozzoErrToMeta