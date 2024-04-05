package svc

import (
	"fmt"
	"time"
)

type ForgotPassSigning struct {
	signingSvc ISigning
}

func NewForgotPassSigningSvc(signingSvc ISigning) ForgotPassSigning {
	return ForgotPassSigning{
		signingSvc: signingSvc,
	}
}

func (s *ForgotPassSigning) Sign(email string) (Token, error) {
	return s.signingSvc.Sign(map[string]any{
		"email": email,
	}, time.Now().Add(time.Hour))
}

func (s *ForgotPassSigning) Parse(token string) (string, error) {
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
