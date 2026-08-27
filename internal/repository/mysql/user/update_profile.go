package user

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/pkg/richerror"
)

func (r Repository) UpdateProfile(ctx context.Context, userID uint, name *string, avatar *string) error {
	const op = "user-repository.UpdateProfile"

	const query = `UPDATE user SET name = COALESCE(?, name), avatar = COALESCE(?, avatar) WHERE id = ?`

	_, err := r.connection.DB.ExecContext(ctx, query, name, avatar, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return richerror.New().
				SetOp(op).
				SetMsg("not found user").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error in update user record").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return nil
}
