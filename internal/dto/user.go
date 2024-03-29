package dto

import "fmt"

type User struct {
	ID                 int32
	Email              string
	Role               string
	expensiveFieldCalc func() int32
}

func NewUser(id int32, email, role string, expensiveFieldCalc func() int32) User {
	fmt.Println("creating user..")
	return User{
		ID:                 id,
		Email:              email,
		Role:               role,
		expensiveFieldCalc: expensiveFieldCalc,
	}
}

func (u User) ExpensiveField() int32 {
	return u.expensiveFieldCalc()
}
