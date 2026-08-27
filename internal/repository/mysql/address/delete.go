package address

import (
	"context"
	"shop/internal/pkg/richerror"
	"time"
)

func (r Repository) Delete(ctx context.Context, userID uint, addressID uint) error {
	const op = "address-repository.Delete"

	const query = `UPDATE address SET deleted_at = ? WHERE user_id = ? AND id = ? AND deleted_at IS NULL`

	_, err := r.connection.DB.ExecContext(ctx, query, time.Now(), userID, addressID)
	if err != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return nil
}
