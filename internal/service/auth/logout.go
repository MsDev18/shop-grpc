package auth

import "context"

func (s Service) Logout (ctx context.Context, sessionID uint) (error) {
	return s.repository.RevokeSession(ctx, sessionID)
}