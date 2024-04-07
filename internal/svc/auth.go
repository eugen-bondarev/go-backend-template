package svc

import (
	"errors"
	"go-backend-template/internal/model"
)

// var (
// 	ErrAuthAuthFailed = util.NewAPIErrorStr(
// 		403,
// 		"authentication failed, please check your credentials",
// 	)
// 	ErrAuthCreateUserFailed = util.NewAPIErrorStr(
// 		403,
// 		"registration failed",
// 	)
// )

var (
	ErrAuthAuthFailed       = errors.New("authentication failed, please check your credentials")
	ErrAuthCreateUserFailed = errors.New("registration failed")
)

type IAuth interface {
	CreateUser(email, plainTextPassword, role string) error
	SetPasswordByEmail(email, plainTextPassword string) error
	AuthenticateUser(email, plainTextPassword string) (model.User, error)
}
