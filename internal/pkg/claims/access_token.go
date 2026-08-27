package claims

import (
	"shop/internal/pkg/richerror"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	SessionID uint `json:"session-id"`
	jwt.RegisteredClaims
}

func CreateAccessToken(userID uint, sessionID uint, secretKey string, duration time.Duration) (string, error) {
	const op = "claims.CreateAccessToken"

	expiresAt := jwt.NewNumericDate(time.Now().Add(duration))
	claims := AccessTokenClaims{
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			Subject:   strconv.Itoa(int(userID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(secretKey)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("signed jwt access token err").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return tokenStr, nil
}

func ParseAccessToken(tokenStr string, secretKey string) (*AccessTokenClaims, error) {
	const op = "claims.ParseAccessToken"
	claims := &AccessTokenClaims{}

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
