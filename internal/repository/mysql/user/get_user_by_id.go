package user

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetUserByID(ctx context.Context, userID uint) (entity.User, error) {
	const op = "user-repository.GetUserByID"

	const query = `SELECT * FROM user WHERE id = ?`

	row := r.connection.DB.QueryRowContext(ctx, query, userID)

	var user entity.User
	err := row.Scan(&user.ID, &user.Role, &user.Name, &user.Avatar, &user.PhoneNumber, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, richerror.New().
				SetOp(op).
				SetMsg("not found user").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.User{}, richerror.New().
			SetOp(op).
			SetMsg("can't scan data from user repository").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return user, nil
}
