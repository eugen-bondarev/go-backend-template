package impl

import (
	"errors"
	"fmt"
	"go-backend-template/internal/model"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTSigningSvc struct {
	secret string
}

func NewJWTSigningSvc(secret string) model.SigningSvc {
	return &JWTSigningSvc{
		secret: secret,
	}
}

func (signingSvc *JWTSigningSvc) Sign(data map[string]any, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  expiration.Unix(),
		"data": data,
	})
	return token.SignedString([]byte(signingSvc.secret))
}

func (signingSvc *JWTSigningSvc) Parse(tokenString string) (map[string]any, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(signingSvc.secret), nil
	})

	if err != nil {
		return map[string]any{}, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if data, ok := claims["data"].(map[string]any); ok {
			return data, nil
		}
	}

	return map[string]any{}, errors.New("failed to validate token")
}
