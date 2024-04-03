package svc

import (
	"errors"
	"time"
)

type UserDataSigningSvc struct {
	signingSvc             ISigningSvc
	sessionTokenExpiration time.Duration
	refreshTokenExpiration time.Duration
}

type RefreshData struct {
	ID        int
	ExpiresAt time.Time
}

type SessionData struct {
	ID        int
	Role      string
	ExpiresAt time.Time
}

func NewUserDataSigningSvc(signingSvc ISigningSvc) UserDataSigningSvc {
	return UserDataSigningSvc{
		signingSvc:             signingSvc,
		sessionTokenExpiration: time.Minute * 2,
		refreshTokenExpiration: time.Minute * 30,
	}
}

func (s *UserDataSigningSvc) SignSessionToken(ID int, role string) (Token, error) {
	return s.signingSvc.Sign(
		map[string]any{
			"ID":   ID,
			"role": role,
		},
		time.Now().Add(s.sessionTokenExpiration),
	)
}

func (s *UserDataSigningSvc) SignRefreshToken(ID int) (Token, error) {
	return s.signingSvc.Sign(
		map[string]any{
			"ID": ID,
		},
		time.Now().Add(s.refreshTokenExpiration),
	)
}

func (s *UserDataSigningSvc) ParseSessionToken(token string) (SessionData, error) {
	data, err := s.signingSvc.Parse(token)

	if err != nil {
		return SessionData{ID: -1}, err
	}

	ID, ok := data["ID"].(float64)

	if !ok {
		return SessionData{ID: -1}, errors.New("token has expired")
	}

	role, ok := data["role"].(string)

	if !ok {
		return SessionData{ID: -1}, errors.New("token has expired")
	}

	exp, ok := data["exp"].(float64)

	if !ok {
		return SessionData{ID: -1}, errors.New("failed to parse exp")
	}

	return SessionData{
		ID:        int(ID),
		Role:      role,
		ExpiresAt: time.Unix(0, int64(exp)*int64(time.Second)),
	}, nil
}

func (s *UserDataSigningSvc) ParseRefreshToken(token string) (RefreshData, error) {
	data, err := s.signingSvc.Parse(token)

	if err != nil {
		return RefreshData{
			ID: -1,
		}, err
	}

	ID := data["ID"].(float64)

	exp, ok := data["exp"].(float64)

	if !ok {
		return RefreshData{}, errors.New("failed to parse exp")
	}

	return RefreshData{
		ID:        int(ID),
		ExpiresAt: time.Unix(0, int64(exp)*int64(time.Second)),
	}, nil
}
