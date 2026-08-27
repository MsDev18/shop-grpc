package province

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetAll(ctx context.Context) ([]entity.Province, error) {
	const op = "province-repository.GetAll"

	const query = `SELECT * FROM province`

	rows, err := r.connection.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer rows.Close()

	var provinceList = make([]entity.Province, 0)
	for rows.Next() {
		var p entity.Province
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, richerror.New().
				SetOp(op).
				SetMsg("can't scan data").
				SetKind(richerror.KindUnexpectedErr).
				SetErr(err)
		}
		provinceList = append(provinceList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return provinceList, nil
}
