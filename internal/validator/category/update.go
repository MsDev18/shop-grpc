package category

import (
	"context"
	"errors"
	dto "shop/internal/dto/category"
	"shop/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) Update(ctx context.Context, req dto.UpdateRequest) error {
	const op = "category-validator.Update"

	err := validation.ValidateStructWithContext(
		ctx,
		&req,
		validation.Field(&req.Title, validation.Length(TITLE_MIN_LENGTH, TITLE_MAX_LENGTH)),
		validation.Field(&req.Slug, validation.Length(SLUG_MIN_LENGTH, SLUG_MAX_LENGTH), validation.Match(SLUG_REGEX)),
		validation.Field(&req.Image, validation.By(v.validateImage)),
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
				SetMsg("invalid input").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err).SetMeta(meta)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindBadRequestErr).
			SetErr(err)
	}
	return nil
}
