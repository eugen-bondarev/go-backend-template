package util

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type APIError struct {
	StatusCode     int
	LocalizeConfig i18n.LocalizeConfig
}

func NewAPIError(statusCode int, localizeConfig i18n.LocalizeConfig) *APIError {
	return &APIError{
		StatusCode:     statusCode,
		LocalizeConfig: localizeConfig,
	}
}

func NewAPIErrorStr(statusCode int, str string) *APIError {
	return NewAPIError(statusCode, i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: str,
		},
	})
}

func (r *APIError) Error() string {
	return r.LocalizeConfig.MessageID
}
