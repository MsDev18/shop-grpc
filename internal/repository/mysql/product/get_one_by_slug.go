package product

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetOneBySlug(ctx context.Context, slug string) (entity.Product, []entity.ProductImage, error) {
	const op = "product-repository.GetBySlug"

	const productQuery = `SELECT * FROM product WHERE slug = ? AND deleted_at IS NULL`

	row := r.connection.DB.QueryRowContext(ctx, productQuery, slug)

	var p entity.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock, &p.MainImage, &p.CategoryID, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, nil, richerror.New().
				SetOp(op).
				SetMsg("not found product").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.Product{}, nil, richerror.New().
			SetOp(op).
			SetMsg("can't scan product data").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	images , err := r.GetProductImage(ctx, p.ID)
	if err != nil {
		return entity.Product{} , nil , err
	}

	return p, images, nil
}
