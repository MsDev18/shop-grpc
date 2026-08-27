package product

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) IsExistsSlug(ctx context.Context, slug string) (bool, error) {
	const op = "product-repository.IsExistsSlug"

	const query = `SELECT * FROM product WHERE slug = ? AND deleted_at IS NULL`

	row := r.connection.DB.QueryRowContext(ctx, query, slug)

	var p entity.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock, &p.MainImage, &p.CategoryID, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	
	return true, nil
}
