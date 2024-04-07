package localization

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18nLocalizer struct {
	bundle     *i18n.Bundle
	localizers map[string]*i18n.Localizer
}

func NewI18nLocalizer() Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("de.toml")
	localizers := make(map[string]*i18n.Localizer)
	localizers["de"] = i18n.NewLocalizer(bundle, "de")
	localizers["en"] = i18n.NewLocalizer(bundle, "en")
	return &I18nLocalizer{
		bundle:     bundle,
		localizers: localizers,
	}
}

func (l *I18nLocalizer) getLocalizer(lang string) *i18n.Localizer {
	res, exists := l.localizers[lang]
	if !exists {
		return l.localizers["en"]
	}
	return res
}

func (l *I18nLocalizer) Translate(msg Message, lang string) string {
	localizer := l.getLocalizer(lang)
	res, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msg.GetID(),
			Other: msg.GetContent(),
		},
	})
	if err != nil {
		return msg.GetContent()
	}
	return res
}
