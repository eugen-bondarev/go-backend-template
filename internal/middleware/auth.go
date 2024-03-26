package middleware

import (
	"fmt"
	"go-backend-template/internal/model"
	"go-backend-template/internal/util"
	"strings"
)

type AuthSuccessHandler = func(ID int, role string)

func Auth(signingSvc model.SigningSvc, authHeader string, successHandler AuthSuccessHandler) error {
	components := strings.Split(authHeader, " ")

	if strings.ToLower(components[0]) != "bearer" {
		return &util.RequestError{
			StatusCode: 403,
			Err:        fmt.Errorf("unauthorized"),
		}
	}

	ID, role, err := signingSvc.Parse(components[1])

	if err != nil {
		return &util.RequestError{
			StatusCode: 403,
			Err:        fmt.Errorf("unauthorized"),
		}
	}

	successHandler(ID, role)

	return nil
}
