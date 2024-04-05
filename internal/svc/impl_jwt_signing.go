package svc

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTSigning struct {
	secret string
}

func NewJWTSigning(secret string) ISigning {
	return &JWTSigning{
		secret: secret,
	}
}

func (signing *JWTSigning) Sign(data map[string]any, expiration time.Time) (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  expiration.Unix(),
		"data": data,
	})

	tokenStr, err := token.SignedString([]byte(signing.secret))
	if err != nil {
		return Token{}, err
	}

	return Token{
		Value:     tokenStr,
		ExpiresAt: expiration,
	}, nil
}

func (signing *JWTSigning) Parse(tokenString string) (map[string]any, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(signing.secret), nil
	})

	if err != nil {
		return map[string]any{}, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if data, ok := claims["data"].(map[string]any); ok {
			data["exp"] = claims["exp"]
			return data, nil
		}
	}

	return map[string]any{}, errors.New("failed to validate token")
}
