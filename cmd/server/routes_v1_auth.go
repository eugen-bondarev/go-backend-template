package main

type RegisterResponse struct {
}

func (app *App) register(email, plainTextPassword string) error {
	return app.authSvc.CreateUser(email, plainTextPassword, "user")
}

func (app *App) login(email, plainTextPassword string) (string, string, error) {
	user, err := app.authSvc.AuthenticateUser(email, plainTextPassword)

	if err != nil {
		return "", "", err
	}

	token, err := app.userDataSigningSvc.SignSessionToken(user.ID, user.Role)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := app.userDataSigningSvc.SignRefreshToken(user.ID)

	if err != nil {
		return "", "", err
	}

	return token.Value, refreshToken.Value, nil
}

func (app *App) refreshToken(refreshToken string) (string, string, error) {
	refreshData, err := app.userDataSigningSvc.ParseRefreshToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	user, err := app.userRepo.GetUserByID(refreshData.ID)

	if err != nil {
		return "", "", err
	}

	token, err := app.userDataSigningSvc.SignSessionToken(refreshData.ID, user.Role)

	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := app.userDataSigningSvc.SignRefreshToken(refreshData.ID)

	if err != nil {
		return "", "", err
	}

	return token.Value, newRefreshToken.Value, nil
}
