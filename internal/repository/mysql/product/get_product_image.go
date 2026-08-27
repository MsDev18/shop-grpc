package product

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetProductImage (ctx context.Context, productID uint) ([]entity.ProductImage, error) {
	const op = "product-repository.GetProductImage"

	const query = `SELECT * FROM product_image WHERE product_id = ? AND deleted_at IS NULL`

	rows, err := r.connection.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer rows.Close()

	var images = make([]entity.ProductImage, 0)
	for rows.Next() {
		var i entity.ProductImage
		if err := rows.Scan(&i.ID, &i.Image, &i.ProductID, &i.DeletedAt, &i.CreatedAt); err != nil {
			return nil, richerror.New().
				SetOp(op).
				SetMsg("can't scan data").
				SetKind(richerror.KindUnexpectedErr).
				SetErr(err)
		}
		images = append(images, i)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return images, nil
}