package main

import "go-backend-template/internal/svc"

type registerResponse = any

func (c *Controller) register(email, plainTextPassword string) (registerResponse, error) {
	err := c.app.auth.CreateUser(email, plainTextPassword, "user")
	return nil, err
}

type loginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (c *Controller) login(email, plainTextPassword string) (loginResponse, error) {
	user, err := c.app.auth.AuthenticateUser(email, plainTextPassword)

	empty := loginResponse{}

	if err != nil {
		return empty, err
	}

	token, err := c.app.userDataSigning.SignSessionToken(user.ID, user.Role)

	if err != nil {
		return empty, err
	}

	refreshToken, err := c.app.userDataSigning.SignRefreshToken(user.ID)

	if err != nil {
		return empty, err
	}

	return loginResponse{
		Token:        token.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}

type refreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (c *Controller) refreshToken(refreshToken string) (refreshTokenResponse, error) {
	refreshData, err := c.app.userDataSigning.ParseRefreshToken(refreshToken)

	empty := refreshTokenResponse{}

	if err != nil {
		return empty, err
	}

	user, err := c.app.userRepo.GetUserByID(refreshData.ID)

	if err != nil {
		return empty, err
	}

	token, err := c.app.userDataSigning.SignSessionToken(refreshData.ID, user.Role)

	if err != nil {
		return empty, err
	}

	newRefreshToken, err := c.app.userDataSigning.SignRefreshToken(refreshData.ID)

	if err != nil {
		return empty, err
	}

	return refreshTokenResponse{
		Token:        token.Value,
		RefreshToken: newRefreshToken.Value,
	}, nil
}

type logoutResponse = any

func (c *Controller) logout(token, refreshToken string) (logoutResponse, error) {
	parsedSessionToken, err := c.app.userDataSigning.ParseSessionToken(token)

	if err != nil {
		return nil, err
	}

	parsedRefreshToken, err := c.app.userDataSigning.ParseRefreshToken(refreshToken)

	if err != nil {
		return nil, err
	}

	c.app.userDataSigning.InvalidateToken(token, parsedSessionToken.ExpiresAt)
	c.app.userDataSigning.InvalidateToken(refreshToken, parsedRefreshToken.ExpiresAt)

	return nil, nil
}

type resetPasswordResponse = any

func (c *Controller) resetPassword(token, password string) (resetPasswordResponse, error) {
	email, err := c.app.forgotPassSigning.Parse(token)

	if err != nil {
		return nil, err
	}

	return nil, c.app.auth.SetPasswordByEmail(email, password)
}

type forgotPasswordResponse = any

func (c *Controller) forgotPassword(email string) (forgotPasswordResponse, error) {
	token, err := c.app.forgotPassSigning.Sign(email)

	if err != nil {
		return nil, err
	}

	mail := svc.NewMailBuilder(
		email,
		"So you want to reset your password?\n"+
			"Your token is: "+token.Value,
	)

	return nil, c.app.mailer.Send(mail)
}
