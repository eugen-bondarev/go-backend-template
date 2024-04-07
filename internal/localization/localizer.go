package localization

type Localizer interface {
	Translate(msg Message, lang string) string
}
