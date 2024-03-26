package middleware

type CanProceedHandler func(role string) bool

type Middleware interface {
	SetRole()
	GetRole() string
	Abort()
}

func AllowIf(m Middleware, handler CanProceedHandler) {
	if !handler(m.GetRole()) {
		m.Abort()
	}
}
