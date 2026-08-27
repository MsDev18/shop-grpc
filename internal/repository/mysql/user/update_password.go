package user

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) UpdatePassword(ctx context.Context, user entity.User) error {
	const op = "user-repository.UpdatePassword"

	const query = "UPDATE user SET password = ? WHERE id = ?"

	_, err := r.connection.DB.ExecContext(ctx, query, user.Password, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return richerror.New().
				SetOp(op).
				SetMsg("not found user with this user id").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return nil
}
