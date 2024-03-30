package model

import "fmt"

type ForgotPassSigningSvc struct {
	signingSvc SigningSvc
}

func NewForgotPassSigningSvc(signingSvc SigningSvc) ForgotPassSigningSvc {
	return ForgotPassSigningSvc{
		signingSvc: signingSvc,
	}
}

func (s *ForgotPassSigningSvc) Sign(email string) (string, error) {
	return s.signingSvc.Sign(map[string]any{
		"email": email,
	})
}

func (s *ForgotPassSigningSvc) Parse(token string) (string, error) {
	data, err := s.signingSvc.Parse(token)

	if err != nil {
		return "", err
	}

	email, ok := data["email"].(string)

	if !ok {
		return "", fmt.Errorf("failed to parse email")
	}

	return email, nil
}
