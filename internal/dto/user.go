package dto

import "fmt"

type User struct {
	ID                 int32
	Email              string
	Role               string
	FirstName          string
	LastName           string
	expensiveFieldCalc func() int32
}

func NewUser(id int32, email, role, firstName, lastName string, expensiveFieldCalc func() int32) User {
	fmt.Println("creating user..")
	return User{
		ID:                 id,
		Email:              email,
		Role:               role,
		FirstName:          firstName,
		LastName:           lastName,
		expensiveFieldCalc: expensiveFieldCalc,
	}
}

func (u User) ExpensiveField() int32 {
	return u.expensiveFieldCalc()
}
