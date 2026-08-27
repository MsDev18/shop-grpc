package address

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) Create(ctx context.Context, address entity.Address) (entity.Address, error) {
	const op = "address-repository.Create"

	const query = `INSERT INTO address (user_id , title , province_id , city , address, postal_code) VALUES (?,?,?,?,?,?)`

	res, err := r.connection.DB.ExecContext(ctx, query, address.UserID, address.Title, address.ProvinceID, address.City, address.Address, address.PostalCode)
	if err != nil {
		return entity.Address{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error in exec query").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	id , _:= res.LastInsertId()
	address.ID = uint(id)
	
	return address, nil
}
