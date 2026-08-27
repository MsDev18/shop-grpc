package auth

import (
	"context"
	"database/sql"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetOtpByUserID(ctx context.Context, userID uint) (entity.Otp, error) {
	const op = "auth-repository.GetOtpByUserID"

	query := `SELECT * FROM otp WHERE user_id = ?`

	var otp entity.Otp

	row := r.connection.DB.QueryRowContext(ctx, query, userID)
	if err := row.Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return entity.Otp{}, richerror.New().
				SetOp(op).
				SetMsg("otp not found").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.Otp{}, richerror.New().
			SetOp(op).
			SetMsg("can't scan otp data").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return otp, nil
}
