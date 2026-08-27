package address

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetOne(ctx context.Context, userID uint, addressID uint) (entity.Address, error) {
	const op = "address-repository.GetOne"

	const query = `SELECT * FROM address WHERE user_id = ? AND id = ? AND deleted_at IS NULL`

	row := r.connection.DB.QueryRowContext(ctx, query, userID, addressID)

	var a entity.Address
	if err := row.Scan(&a.ID, &a.UserID, &a.Title, &a.ProvinceID, &a.City, &a.Address, &a.PostalCode, &a.DeletedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Address{}, richerror.New().
				SetOp(op).
				SetMsg("not found address with this id").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.Address{}, richerror.New().
			SetOp(op).
			SetMsg("can't scan address data").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return a, nil
}
