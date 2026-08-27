package auth

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (r Repository) CreateSession(ctx context.Context, session entity.Session) (entity.Session, error) {
	const op = "auth-repository.CreateSession"

	const query = `INSERT INTO session (user_id,expires_at) VALUES (?,?)`
	result, err := r.connection.DB.ExecContext(ctx, query, session.UserID, session.ExpiresAt)
	if err != nil {
		return entity.Session{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error in create session in database").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	sessionID, _ := result.LastInsertId()
	session.ID = uint(sessionID)

	return session, nil
}
