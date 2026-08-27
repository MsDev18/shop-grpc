package address

import (
	"context"
	"shop/internal/pkg/richerror"
)

func (r Repository) Update(ctx context.Context, userID uint, addressID uint, title *string, provinceID *uint, city, address, postalCode *string) error {
	const op = "address-repository.Update"

	const query = `
		UPDATE address SET 
			title = COALESCE (? ,title), 
			province_id = COALESCE(?, province_id),
			city = COALESCE(? , city), 
			address = COALESCE(? , address), 
			postal_code = COALESCE(? , postal_code) 
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL 
	`

	_, err := r.connection.DB.ExecContext(ctx, query, title, provinceID, city, address, postalCode, addressID, userID)
	if err != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return nil
}
