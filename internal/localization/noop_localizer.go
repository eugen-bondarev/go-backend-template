package localization

type NoopLocalizer struct {
}

func NewNoopLocalizer() Localizer {
	return &NoopLocalizer{}
}

func (l *NoopLocalizer) Translate(msg Message, lang string) string {
	return msg.content
}
