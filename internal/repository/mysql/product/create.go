package product

import (
	"context"
	"database/sql"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) Create(ctx context.Context, p entity.Product, imagePaths []string) (entity.Product, error) {
	const op = "product-repository.Create"

	tx, err := r.connection.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.Product{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer tx.Rollback()

	product, err := r.insertProduct(ctx, tx, p)
	if err != nil {
		return entity.Product{}, err
	}

	if err := r.insertProductImages(ctx, tx, product.ID, imagePaths); err != nil {
		return entity.Product{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.Product{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return product, nil
}

func (r Repository) insertProduct(ctx context.Context, tx *sql.Tx, p entity.Product) (entity.Product, error) {
	const op = "product-repository.insertProduct"

	const productQuery = `
		INSERT INTO product (name , slug , description, price , stock , main_image, category_id) VALUES (?,?,?,?,?,?,?)
	`

	res, err := tx.ExecContext(ctx, productQuery, p.Name, p.Slug, p.Description, p.Price, p.Stock, p.MainImage, p.CategoryID)
	if err != nil {
		return entity.Product{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return entity.Product{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	p.ID = uint(id)

	return p, nil
}

func (r Repository) insertProductImages(ctx context.Context, tx *sql.Tx, productID uint, imagePaths []string) error {
	const op = "product-repository.insertProductImages"

	const imageQuery = `INSERT INTO product_image (product_id , image) VALUES (? ,?)`

	for _, imagePath := range imagePaths {
		if _, err := tx.ExecContext(ctx, imageQuery, productID, imagePath); err != nil {
			return richerror.New().
				SetOp(op).
				SetMsg("unexpected error").
				SetKind(richerror.KindUnexpectedErr).
				SetErr(err)
		}
	}

	return nil
}
