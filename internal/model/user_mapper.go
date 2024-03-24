package model

type UserMapper[T any] interface {
	FromUser(User) T
	ToUser(T) User
}
