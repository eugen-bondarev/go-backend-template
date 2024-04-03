package middleware

import (
	"go-backend-template/internal/permissions"
	"go-backend-template/internal/svc"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinMiddleware struct {
	ctx                 *gin.Context
	userDataSigningSvc  *svc.UserDataSigningSvc
	tokenInvalidatorSvc svc.ITokenInvalidatorSvc
}

func NewGinMiddleware(
	ctx *gin.Context,
	userDataSigningSvc *svc.UserDataSigningSvc,
	tokenInvalidatorSvc svc.ITokenInvalidatorSvc,
) Middleware {
	return &GinMiddleware{
		ctx:                 ctx,
		userDataSigningSvc:  userDataSigningSvc,
		tokenInvalidatorSvc: tokenInvalidatorSvc,
	}
}

func (m *GinMiddleware) getRoleFromParams() string {
	return m.ctx.Query("foo")
}

func (m *GinMiddleware) getRoleFromHeader() string {
	authHeader := m.ctx.Request.Header.Get("Authorization")
	components := strings.Split(authHeader, " ")

	if strings.ToLower(components[0]) != "bearer" {
		return ""
	}

	sessionData, _ := m.userDataSigningSvc.ParseSessionToken(components[1])

	if !m.tokenInvalidatorSvc.IsValid(components[1]) {
		return ""
	}

	return sessionData.Role
}

func (m *GinMiddleware) SetRole() {
	m.ctx.Set("role", m.getRoleFromHeader())
}

func (m *GinMiddleware) GetRole() string {
	return m.ctx.GetString("role")
}

func (m *GinMiddleware) Abort() {
	m.ctx.Abort()
	m.ctx.JSON(403, gin.H{
		"error": "unauthorized",
	})
}

type GinMiddlewareFactory struct {
	userDataSigningSvc svc.UserDataSigningSvc
	tokenInvalidator   svc.ITokenInvalidatorSvc
	policies           *permissions.Policies
}

func NewGinMiddlewareFactory(
	userDataSigningSvc svc.UserDataSigningSvc,
	tokenInvalidator svc.ITokenInvalidatorSvc,
	policies *permissions.Policies,
) GinMiddlewareFactory {
	return GinMiddlewareFactory{
		userDataSigningSvc: userDataSigningSvc,
		tokenInvalidator:   tokenInvalidator,
		policies:           policies,
	}
}

func (factory *GinMiddlewareFactory) SetRole() func(*gin.Context) {
	return func(ctx *gin.Context) {
		m := NewGinMiddleware(ctx, &factory.userDataSigningSvc, factory.tokenInvalidator)
		m.SetRole()
	}
}

func (factory *GinMiddlewareFactory) EnforcePolicy(action, object string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		m := NewGinMiddleware(ctx, nil, factory.tokenInvalidator)
		AllowIf(m, func(role string) bool {
			return factory.policies.RoleCan(role, action, object)
		})
	}
}
