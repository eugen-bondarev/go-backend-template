package model

import "time"

type UserDataSigningSvc struct {
	signingSvc SigningSvc
}

func NewUserDataSigningSvc(signingSvc SigningSvc) UserDataSigningSvc {
	return UserDataSigningSvc{
		signingSvc: signingSvc,
	}
}

func (s *UserDataSigningSvc) Sign(ID int, role string) (string, error) {
	return s.signingSvc.Sign(map[string]any{
		"ID":   ID,
		"role": role,
	}, time.Now().Add(time.Hour))
}

func (s *UserDataSigningSvc) Parse(token string) (int, string, error) {
	data, err := s.signingSvc.Parse(token)

	if err != nil {
		return -1, "", err
	}

	ID := data["ID"].(float64)
	role := data["role"].(string)

	return int(ID), role, nil
}
