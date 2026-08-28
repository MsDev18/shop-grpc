package auth

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/claims"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const (
	USER_ID_KEY    ctxKey = "user-id"
	SESSION_ID_KEY ctxKey = "session-id"
	ROLE_KEY       ctxKey = "role"
)

func UserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(USER_ID_KEY).(uint)
	return userID, ok
}

func SessionIDFromContext(ctx context.Context) (uint, bool) {
	sessionID, ok := ctx.Value(SESSION_ID_KEY).(uint)
	return sessionID, ok
}

func RoleFromContext(ctx context.Context) (entity.Role, bool) {
	role, ok := ctx.Value(ROLE_KEY).(entity.Role)
	return role, ok
}

type Repository interface {
	GetSessionByID(ctx context.Context, sessionID uint) (entity.Session, error)
	GetUserByID(ctx context.Context, userID uint) (entity.User, error)
}

type Interceptor struct {
	repository        Repository
	accessTokenSecret string
}

func NewInterceptor(repository Repository, accessTokenSecret string) Interceptor {
	return Interceptor{
		repository:        repository,
		accessTokenSecret: accessTokenSecret,
	}
}

var publicMethods = map[string]bool{
	"/auth/AuthService.SendOtp":      true,
	"/auth/AuthService.CheckOtp":     true,
	"/auth/AuthService.RefreshToken": true,
	"health/HealthService.Check":     true,
}

func (i Interceptor) Unary() grpc.UnaryServerInterceptor {
	const op = "auth-interceptor"

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			mapper.ErrorToGrpc(
				richerror.New().
					SetOp(op).
					SetMsg("login required").
					SetKind(richerror.KindUnauthorizeErr),
			)
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, mapper.ErrorToGrpc(
				richerror.New().
					SetOp(op).
					SetMsg("login required").
					SetKind(richerror.KindUnauthorizeErr),
			)
		}

		parts := strings.SplitN(values[0], " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return nil, mapper.ErrorToGrpc(
				richerror.New().
					SetOp(op).
					SetMsg("invalid authorization header format").
					SetKind(richerror.KindUnauthorizeErr),
			)
		}

		accessClaims, err := claims.ParseAccessToken(parts[1], i.accessTokenSecret)
		if err != nil {
			return nil, mapper.ErrorToGrpc(err)
		}

		userID, err := strconv.ParseUint(accessClaims.Subject, 10, 64)
		if err != nil {
			return nil, mapper.ErrorToGrpc(
				richerror.New().
					SetOp(op).
					SetMsg("invalid subject claim in access token").
					SetKind(richerror.KindUnauthorizeErr).
					SetErr(err),
			)
		}

		session, err := i.repository.GetSessionByID(ctx, accessClaims.SessionID)
		if err != nil {
			return nil, mapper.ErrorToGrpc(err)
		}

		if session.RevokeAt != nil || !session.ExpiresAt.After(time.Now()) {
			return nil, mapper.ErrorToGrpc(
				richerror.New().
					SetOp(op).
					SetMsg("session revoked or expired").
					SetKind(richerror.KindUnauthorizeErr),
			)
		}

		user, err := i.repository.GetUserByID(ctx, session.UserID)
		if err != nil {
			return nil, mapper.ErrorToGrpc(err)
		}

		ctx = context.WithValue(ctx, USER_ID_KEY, uint(userID))
		ctx = context.WithValue(ctx, SESSION_ID_KEY, uint(session.ID))
		ctx = context.WithValue(ctx, ROLE_KEY, user.Role)

		return handler(ctx, req)
	}
}
