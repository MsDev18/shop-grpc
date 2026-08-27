package auth

import (
	"context"
	"shop/internal/pkg/richerror"
	"time"
)

func (r Repository) RevokeSession(ctx context.Context, sessionID uint) error {
	const op = "auth-repository.RevokeSession"

	const query = `UPDATE session SET revoke_at = ? WHERE id = ?`

	now := time.Now()
	if _, err := r.connection.DB.ExecContext(ctx, query, now, sessionID); err != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error in revoke session in database").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return nil
}
