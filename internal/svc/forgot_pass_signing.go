package svc

import (
	"fmt"
	"time"
)

type ForgotPassSigning struct {
	signing ISigning
}

func NewForgotPassSigning(signing ISigning) ForgotPassSigning {
	return ForgotPassSigning{
		signing: signing,
	}
}

func (s *ForgotPassSigning) Sign(email string) (Token, error) {
	return s.signing.Sign(map[string]any{
		"email": email,
	}, time.Now().Add(time.Hour))
}

func (s *ForgotPassSigning) Parse(token string) (string, error) {
	data, err := s.signing.Parse(token)

	if err != nil {
		return "", err
	}

	email, ok := data["email"].(string)

	if !ok {
		return "", fmt.Errorf("failed to parse email")
	}

	return email, nil
}
