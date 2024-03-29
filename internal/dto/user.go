package dto

type User struct {
	ID        int32
	Email     string
	Role      string
	FirstName string
	LastName  string
}

func NewUser(id int32, email, role, firstName, lastName string) User {
	return User{
		ID:        id,
		Email:     email,
		Role:      role,
		FirstName: firstName,
		LastName:  lastName,
	}
}
