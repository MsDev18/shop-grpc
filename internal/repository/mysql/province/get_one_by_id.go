package province

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetOneByID(ctx context.Context, id uint) (entity.Province, error) {
	const op = "province-repository.GetOneByID"

	const query = `SELECT * FROM province WHERE id = ?`

	row := r.connection.DB.QueryRowContext(ctx, query, id)

	var province entity.Province
	if err := row.Scan(&province.ID, &province.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Province{}, richerror.New().
				SetOp(op).
				SetMsg("not found province with this id").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.Province{}, richerror.New().
			SetOp(op).
			SetMsg("can't scan data").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return province, nil
}
