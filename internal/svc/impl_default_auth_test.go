package svc

import (
	"go-backend-template/internal/repo"
	"testing"
)

type testUser struct {
	email             string
	plainTextPassword string
	role              string
}

type authModule struct {
	auth     IAuth
	userRepo repo.IUserRepo
	testUser testUser
}

func newAuthModule() authModule {
	userRepo := repo.NewMemUserRepo()
	auth := NewDefaultAuth(userRepo, "have you heard what they said on the news today?")
	return authModule{
		auth:     auth,
		userRepo: userRepo,
		testUser: testUser{
			email:             "admin@example.com",
			plainTextPassword: "lorem ipsum",
			role:              "admin",
		},
	}
}

func Test_CanFailAuth(t *testing.T) {
	authModule := newAuthModule()

	_, err := authModule.auth.AuthenticateUser(
		authModule.testUser.email,
		authModule.testUser.plainTextPassword,
	)
	if err == nil {
		t.Errorf("expected err, got nil")
	}

	_, err = authModule.auth.AuthenticateUser(authModule.testUser.email, "foobar")
	if err == nil {
		t.Errorf("expected err, got nil")
	}
}

func Test_CanCreateUser(t *testing.T) {
	authModule := newAuthModule()

	err := authModule.auth.CreateUser(
		authModule.testUser.email,
		authModule.testUser.plainTextPassword,
		authModule.testUser.role,
	)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_CanAuthenticateUser(t *testing.T) {
	authModule := newAuthModule()

	err := authModule.auth.CreateUser(
		authModule.testUser.email,
		authModule.testUser.plainTextPassword,
		authModule.testUser.role,
	)
	if err != nil {
		t.Errorf(err.Error())
	}

	authenticatedUser, err := authModule.auth.AuthenticateUser(
		authModule.testUser.email,
		authModule.testUser.plainTextPassword,
	)
	if err != nil {
		t.Errorf(err.Error())
	}

	if authenticatedUser.Role != authModule.testUser.role {
		t.Errorf("expected %v, got %v", authModule.testUser.role, authenticatedUser.Role)
	}
}

func Test_PasswordIsNotPlainText(t *testing.T) {
	authModule := newAuthModule()

	err := authModule.auth.CreateUser(
		authModule.testUser.email,
		authModule.testUser.plainTextPassword,
		authModule.testUser.role,
	)

	if err != nil {
		t.Errorf(err.Error())
	}

	user, err := authModule.userRepo.GetUserByEmail(authModule.testUser.email)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(user.PasswordHash) == 0 {
		t.Errorf("password is empty")
	}

	if user.PasswordHash == authModule.testUser.plainTextPassword {
		t.Errorf("password matches hash: %v", user.PasswordHash)
	}
}
