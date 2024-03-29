package dto

type User struct {
	ID    int32
	Email string
	Role  string
}

func NewUser(id int32, email, role string) User {
	return User{
		ID:    id,
		Email: email,
		Role:  role,
	}
}
