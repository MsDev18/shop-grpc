package address

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetAll(ctx context.Context, userID uint) ([]entity.Address, error) {
	const op = "address-repository.GetAll"

	const query = `SELECT * FROM address WHERE user_id = ? AND deleted_at IS NULL`

	rows, err := r.connection.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer rows.Close()

	addressList := make([]entity.Address, 0)
	for rows.Next() {
		var a entity.Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.Title, &a.ProvinceID, &a.City, &a.Address, &a.PostalCode, &a.DeletedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, richerror.New().
				SetOp(op).
				SetMsg("can't scan data").
				SetKind(richerror.KindUnexpectedErr).
				SetErr(err)
		}
		addressList = append(addressList, a)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return addressList , nil
}
