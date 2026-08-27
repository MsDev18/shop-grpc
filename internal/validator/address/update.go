package address

import (
	"context"
	"errors"
	dto "shop/internal/dto/address"
	"shop/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) Update(ctx context.Context, req dto.UpdateRequest) error {
	const op = "address-validator.Update"

	err := validation.ValidateStructWithContext(
		ctx,
		&req,
		validation.Field(&req.Title, validation.Length(3, 200)),
		validation.Field(&req.ProvinceID, validation.Match(PROVINCE_ID_REGEX)),
		validation.Field(&req.City, validation.Length(3, 200)),
		validation.Field(&req.Address, validation.Length(10, 500)),
		validation.Field(&req.PostalCode, validation.Match(POSTAL_CODE_REGEX)),
	)
	if err != nil {
		var validationErr *validation.Errors
		if errors.As(err, &validationErr) {
			meta := make(map[string]any)
			for key, value := range *validationErr {
				meta[key] = value.Error()
			}
			return richerror.New().
				SetOp(op).
				SetMsg("invalid inputs").
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
