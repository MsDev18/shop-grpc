package claims

import (
	"shop/internal/pkg/richerror"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenClaims struct {
	SessionID uint `json:"session-id"`
	jwt.RegisteredClaims
}

func CreateRefreshToken(userID uint,sessionID uint, secretKey string, duration time.Duration) (string, error) {
	const op = "claims.CreateRefreshToken"

	expiresAt := jwt.NewNumericDate(time.Now().Add(duration))
	claims := RefreshTokenClaims{
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			Subject: strconv.Itoa(int(userID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(secretKey)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("signed jwt refresh token err").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return tokenStr, nil
}

func ParseRefreshToken(tokenStr string, secretKey string) (*RefreshTokenClaims, error) {
	const op = "claims.ParseRefreshToken"

	claims := &RefreshTokenClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, richerror.New().
					SetOp(op).
					SetMsg("unexpected signing method").
					SetKind(richerror.KindUnauthorizeErr)
			}
			
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("can't parse jwt token").
			SetKind(richerror.KindUnauthorizeErr).
			SetErr(err)
	}

	if !token.Valid {
		return nil, richerror.New().
			SetOp(op).
			SetMsg("token is not valid").
			SetKind(richerror.KindUnauthorizeErr)
	}

	return claims, nil
}