package main

func (c *Controller) register(email, plainTextPassword string) error {
	return c.app.authSvc.CreateUser(email, plainTextPassword, "user")
}

func (c *Controller) login(email, plainTextPassword string) (string, string, error) {
	user, err := c.app.authSvc.AuthenticateUser(email, plainTextPassword)

	if err != nil {
		return "", "", err
	}

	token, err := c.app.userDataSigningSvc.SignSessionToken(user.ID, user.Role)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := c.app.userDataSigningSvc.SignRefreshToken(user.ID)

	if err != nil {
		return "", "", err
	}

	return token.Value, refreshToken.Value, nil
}

func (c *Controller) refreshToken(refreshToken string) (string, string, error) {
	refreshData, err := c.app.userDataSigningSvc.ParseRefreshToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	user, err := c.app.userRepo.GetUserByID(refreshData.ID)

	if err != nil {
		return "", "", err
	}

	token, err := c.app.userDataSigningSvc.SignSessionToken(refreshData.ID, user.Role)

	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := c.app.userDataSigningSvc.SignRefreshToken(refreshData.ID)

	if err != nil {
		return "", "", err
	}

	return token.Value, newRefreshToken.Value, nil
}
