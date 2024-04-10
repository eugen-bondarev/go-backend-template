package repo

import (
	"testing"
)

func Test_UserRepo(t *testing.T) {
	userRepo := NewMemUserRepo()

	type testUser struct {
		email    string
		password string
		role     string
	}

	testUsers := []testUser{
		{
			email:    "admin@example.com",
			password: "fooooo",
			role:     "admin",
		},
		{
			email:    "john.doe@example.com",
			password: "bar",
			role:     "user",
		},
	}

	for _, user := range testUsers {
		userRepo.CreateUser(user.email, user.password, user.role)
	}

	users, _ := userRepo.GetUsers()
	if len(users) != len(testUsers) {
		t.Errorf("expected %v, got %v", len(testUsers), len(users))
	}

	user, err := userRepo.GetUserByEmail(testUsers[0].email)
	if err != nil {
		t.Errorf(err.Error())
	}

	if user.Email != testUsers[0].email {
		t.Errorf("expected %v, got %v", testUsers[0].email, user.Email)
	}
}
