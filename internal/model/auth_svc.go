package model

import (
	"go-backend-template/internal/util"
)

var (
	ErrAuthSvcAuthFailed = util.NewRequestErrorStr(
		403,
		"authentication failed, please check your credentials",
	)
	ErrAuthSvcCreateUserFailed = util.NewRequestErrorStr(
		403,
		"registration failed",
	)
)

type AuthSvc interface {
	CreateUser(email, plainTextPassword, role string) error
	AuthenticateUser(email, plainTextPassword string) (User, error)
}
