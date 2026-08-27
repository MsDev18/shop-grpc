package product

import (
	"context"
	"fmt"
	"strings"

	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetAll(ctx context.Context, limit int, offset int, categoryID *uint) ([]entity.Product, int, error) {
	const op = "product-repository.GetAll"

	conditions := []string{"deleted_at IS NULL"}
	args := []any{}

	if categoryID != nil {
		conditions = append(conditions, "category_id = ?")
		args = append(args, *categoryID)
	}

	whereClause := strings.Join(conditions, " AND ")

	dataQuery := fmt.Sprintf(`SELECT * FROM product WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?`, whereClause)

	dataArgs := make([]any, 0, len(args)+2)
	dataArgs = append(dataArgs, args...)
	dataArgs = append(dataArgs, limit, offset)

	rows, err := r.connection.DB.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer rows.Close()

	var products = make([]entity.Product, 0)
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock, &p.MainImage, &p.CategoryID, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, richerror.New().
				SetOp(op).
				SetMsg("can't scan product data").
				SetKind(richerror.KindUnexpectedErr).
				SetErr(err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, richerror.New().
			SetOp(op).
			SetMsg("error while iterating product rows").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM product WHERE %s`, whereClause)

	var total int
	if err := r.connection.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, richerror.New().
			SetOp(op).
			SetMsg("can't count products").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return products, total, nil
}