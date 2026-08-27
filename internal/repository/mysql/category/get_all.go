package category

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetAll(ctx context.Context) ([]entity.Category, error) {
	const op = "category-repository.GetAll"

	const query = `SELECT * FROM category WHERE deleted_at IS NULL`

	rows, err := r.connection.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer rows.Close()

	var categories = make([]entity.Category,0)
	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Title, &c.Slug, &c.Image, &c.DeletedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, richerror.New().
				SetOp(op).
				SetMsg("can't scan category").
				SetKind(richerror.KindUnexpectedErr).
				SetErr(err)
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("error while iterating category rows").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return categories, nil
}
