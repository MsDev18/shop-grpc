package category

import (
	"context"
	"errors"
	"shop/internal/pkg/richerror"

	"github.com/go-sql-driver/mysql"
)

func (r Repository) UpdateByID(ctx context.Context, id uint, title, slug, image *string) error {
	const op = "category-repository.UpdateByID"

	const query = `
	UPDATE category SET 
		title = COALESCE(?, title) ,
		slug = COALESCE (?,slug), 
		image = COALESCE(?, image)
	WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.connection.DB.ExecContext(ctx, query, title, slug, image, id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return richerror.New().
				SetOp(op).
				SetMsg("duplicate entry").
				SetKind(richerror.KindConflictErr).
				SetErr(err)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return nil
}
