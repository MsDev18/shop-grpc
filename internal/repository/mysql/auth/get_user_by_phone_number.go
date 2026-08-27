package auth

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error) {
	const op = "auth-repository.GetUserByPhoneNumber"

	var user entity.User
	query := `SELECT * FROM user WHERE phone_number = ?`

	row := r.connection.DB.QueryRowContext(ctx, query, phoneNumber)
	err := row.Scan(&user.ID, &user.Role, &user.Name, &user.Avatar, &user.PhoneNumber, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, richerror.New().
				SetOp(op).
				SetMsg("not found user with this phone number").
				SetKind(richerror.KindNotFoundErr)
		}
		return entity.User{}, richerror.New().
			SetOp(op).
			SetMsg("can't scan data from database").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return user, nil
}
