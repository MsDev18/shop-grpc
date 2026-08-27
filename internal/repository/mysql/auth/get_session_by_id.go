package auth

import (
	"context"
	"database/sql"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) GetSessionByID(ctx context.Context, sessionID uint) (entity.Session, error) {
	const op = "auth-repository.GetSessionByID"

	const query = `SELECT * FROM session WHERE id = ?`

	row := r.connection.DB.QueryRowContext(ctx, query, sessionID)

	var session entity.Session
	if err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.RevokeAt); err != nil {
		if err == sql.ErrNoRows {
			return entity.Session{}, richerror.New().
				SetOp(op).
				SetMsg("session not found").
				SetKind(richerror.KindNotFoundErr).
				SetErr(err)
		}
		return entity.Session{}, richerror.New().
			SetOp(op).
			SetMsg("can't scan session data").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return session, nil
}
