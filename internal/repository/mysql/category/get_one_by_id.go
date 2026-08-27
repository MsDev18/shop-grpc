package category

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetOneByID(ctx context.Context, id uint) (entity.Category, error) {
	const op = "category-repository.GetOneByID"

	const query = `SELECT * FROM category WHERE id = ? AND deleted_at IS NULL`

	row := r.connection.DB.QueryRowContext(ctx, query, id)

	var c entity.Category
	if err := row.Scan(&c.ID, &c.ParentID, &c.Title, &c.Slug, &c.Image, &c.DeletedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, richerror.New().
				SetOp(op).
				SetMsg("not found category with this id").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.Category{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return c , nil
}
