package auth

import (
	"shop/internal/pkg/claims"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)


func (m Middleware) AuthRequired() gin.HandlerFunc {
	const op = "auth-middleware"
	return func(ctx *gin.Context) {
		// get bearer token form autorization
		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			response.New(ctx).Error(richerror.New().
				SetOp(op).
				SetMsg("login required").
				SetKind(richerror.KindUnauthorizeErr),
			)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.New(ctx).Error(richerror.New().
				SetOp(op).
				SetMsg("invalid authorization header format").
				SetKind(richerror.KindUnauthorizeErr),
			)
			ctx.Abort()
			return
		}

		token := parts[1]
		accessClaims, err := claims.ParseAccessToken(token, m.accessTokenSecret)
		if err != nil {
			response.New(ctx).Error(err)
			ctx.Abort()
			return
		}

		userID, err := strconv.ParseUint(accessClaims.Subject, 10, 64)
		if err != nil {
			response.New(ctx).Error(richerror.New().
				SetOp(op).
				SetMsg("invalid subject claim in access token").
				SetKind(richerror.KindUnauthorizeErr).
				SetErr(err),
			)
			ctx.Abort()
			return
		}

		session, err := m.repository.GetSessionByID(ctx, accessClaims.SessionID)
		if err != nil {
			response.New(ctx).Error(err)
			ctx.Abort()
			return
		}

		if session.RevokeAt != nil || !session.ExpiresAt.After(time.Now()) {
			response.New(ctx).Error(richerror.New().
				SetOp(op).
				SetMsg("session revoked or expired").
				SetKind(richerror.KindUnauthorizeErr),
			)
			ctx.Abort()
			return
		}

		user, err := m.repository.GetUserByID(ctx, session.UserID)
		if err != nil {
			response.New(ctx).Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(SESSION_ID_KEY, session.ID)
		ctx.Set(USER_ID_KEY, uint(userID))
		ctx.Set(ROLE_KEY, user.Role)
		
		ctx.Next()
	}
}
