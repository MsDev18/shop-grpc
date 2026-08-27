package category

import (
	"context"
	"shop/internal/pkg/richerror"
	"time"
)

func (r Repository) DeleteByID(ctx context.Context, id uint) error {
	const op = "category-repository.DeleteByID"

	const query = `UPDATE category SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`

	_, err := r.connection.DB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	
	return nil
}
