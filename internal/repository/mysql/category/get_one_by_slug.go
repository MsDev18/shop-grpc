package category

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetOneBySlug(ctx context.Context, slug string) (entity.Category, error) {
	const op = "category-repository.GetOneBySlug"

	const query = `SELECT * FROM category WHERE slug = ? AND deleted_at IS NULL`

	row := r.connection.DB.QueryRowContext(ctx, query, slug)

	var c entity.Category
	if err := row.Scan(&c.ID, &c.ParentID, &c.Title, &c.Slug, &c.Image, &c.DeletedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, richerror.New().
				SetOp(op).
				SetMsg("not found category").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.Category{}, richerror.New().
			SetOp(op).
			SetMsg("can't scan category").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	
	return c, nil
}
