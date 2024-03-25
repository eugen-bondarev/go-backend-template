package impl

import (
	"errors"
	"fmt"
	"go-backend-template/internal/model"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
)

const invalidID = -1
const invalidToken = ""

type JWTSigningSvc struct {
	secret string
}

func NewJWTSigningSvc(secret string) model.SigningSvc {
	return &JWTSigningSvc{
		secret: secret,
	}
}

type tokenBody struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

func (signingSvc *JWTSigningSvc) Sign(ID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenBody": tokenBody{
			ID:   ID,
			Role: role,
		},
	})
	return token.SignedString([]byte(signingSvc.secret))
}

func (signingSvc *JWTSigningSvc) Parse(tokenString string) (int, string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(signingSvc.secret), nil
	})

	if err != nil {
		return invalidID, invalidToken, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		tb := struct {
			TokenBody tokenBody `json:"tokenBody"`
		}{}
		mapstructure.Decode(claims, &tb)
		return tb.TokenBody.ID, tb.TokenBody.Role, nil
	}

	return invalidID, invalidToken, errors.New("failed to validate token")
}
