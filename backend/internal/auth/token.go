package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenExpiry = 48 * time.Hour // 2 days
var ErrInvalidToken = errors.New("invalid or expired token")


type claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, secret string) (string, error) {
	now := time.Now()
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}

func ParseToken(tokenString, secret string) (int64, error) {
	c := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	return c.UserID, nil
}

