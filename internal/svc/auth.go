package svc

import (
	"go-backend-template/internal/model"
	"go-backend-template/internal/util"
)

var (
	ErrAuthAuthFailed = util.NewAPIErrorStr(
		403,
		"authentication failed, please check your credentials",
	)
	ErrAuthCreateUserFailed = util.NewAPIErrorStr(
		403,
		"registration failed",
	)
)

type IAuth interface {
	CreateUser(email, plainTextPassword, role string) error
	SetPasswordByEmail(email, plainTextPassword string) error
	AuthenticateUser(email, plainTextPassword string) (model.User, error)
}
