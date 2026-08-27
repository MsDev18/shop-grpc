package auth

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) UpsertOtp(ctx context.Context, otp entity.Otp) (entity.Otp, error) {
	const op = "auth-repository.UpsertOtp"

	const query = `
	INSERT INTO otp (user_id, code, expires_at) 
	VALUES(?,?,?) AS new
	ON DUPLICATE KEY UPDATE
		id = LAST_INSERT_ID(id), 
		code = new.code, 
		expires_at = new.expires_at
	`

	result , resultErr := r.connection.DB.ExecContext(ctx, query, otp.UserID, otp.Code, otp.ExpiresAt)
	if resultErr != nil {
		return entity.Otp{}, richerror.New().
		SetOp(op).
		SetMsg("unexpected error").
		SetKind(richerror.KindUnexpectedErr).
		SetErr(resultErr)
	}
	otpID ,_ := result.LastInsertId()
	otp.ID = uint(otpID)

	return otp, nil
}
