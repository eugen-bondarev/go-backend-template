package main

type RegisterResponse struct {
}

func (app *App) register(email, plainTextPassword string) error {
	return app.authSvc.CreateUser(email, plainTextPassword, "user")
}

func (app *App) login(email, plainTextPassword string) (string, error) {
	user, err := app.authSvc.AuthenticateUser(email, plainTextPassword)

	if err != nil {
		return "", err
	}

	token, err := app.userDataSigningSvc.Sign(user.ID, user.Role)

	return token, err
}
