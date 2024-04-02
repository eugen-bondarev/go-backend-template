package svc

import (
	"go-backend-template/internal/model"
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

type IAuthSvc interface {
	CreateUser(email, plainTextPassword, role string) error
	SetPasswordByEmail(email, plainTextPassword string) error
	AuthenticateUser(email, plainTextPassword string) (model.User, error)
}
