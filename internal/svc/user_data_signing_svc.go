package svc

import (
	"errors"
	"time"
)

type IUserDataSigningSvc struct {
	signingSvc             ISigningSvc
	sessionTokenExpiration time.Duration
	refreshTokenExpiration time.Duration
}

func NewUserDataSigningSvc(signingSvc ISigningSvc) IUserDataSigningSvc {
	return IUserDataSigningSvc{
		signingSvc:             signingSvc,
		sessionTokenExpiration: time.Minute * 2,
		refreshTokenExpiration: time.Minute * 30,
	}
}

func (s *IUserDataSigningSvc) SignSessionToken(ID int, role string) (string, error) {
	token, err := s.signingSvc.Sign(
		map[string]any{
			"ID":   ID,
			"role": role,
		},
		time.Now().Add(s.sessionTokenExpiration),
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *IUserDataSigningSvc) ParseSessionToken(token string) (int, string, error) {
	data, err := s.signingSvc.Parse(token)

	if err != nil {
		return -1, "", err
	}

	ID, ok := data["ID"].(float64)

	if !ok {
		return -1, "", errors.New("token has expired")
	}

	role, ok := data["role"].(string)

	if !ok {
		return -1, "", errors.New("token has expired")
	}

	return int(ID), role, nil
}

func (s *IUserDataSigningSvc) SignRefreshToken(ID int) (string, error) {
	refreshToken, err := s.signingSvc.Sign(
		map[string]any{
			"ID": ID,
		},
		time.Now().Add(s.refreshTokenExpiration),
	)

	return refreshToken, err
}

func (s *IUserDataSigningSvc) ParseRefreshToken(token string) (int, error) {
	data, err := s.signingSvc.Parse(token)

	if err != nil {
		return -1, err
	}

	ID := data["ID"].(float64)

	return int(ID), nil
}
