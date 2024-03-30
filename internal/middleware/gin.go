package middleware

import (
	"fmt"
	"go-backend-template/internal/impl"
	"go-backend-template/internal/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinMiddleware struct {
	ctx                *gin.Context
	userDataSigningSvc *model.UserDataSigningSvc
}

func NewGinMiddleware(ctx *gin.Context, userDataSigningSvc *model.UserDataSigningSvc) Middleware {
	return &GinMiddleware{
		ctx:                ctx,
		userDataSigningSvc: userDataSigningSvc,
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

	_, role, err := m.userDataSigningSvc.Parse(components[1])

	fmt.Println(err)

	return role
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
	userDataSigningSvc model.UserDataSigningSvc
	policies           *impl.Policies
}

func NewGinMiddlewareFactory(userDataSigningSvc model.UserDataSigningSvc, policies *impl.Policies) GinMiddlewareFactory {
	return GinMiddlewareFactory{
		userDataSigningSvc: userDataSigningSvc,
		policies:           policies,
	}
}

func (factory *GinMiddlewareFactory) SetRole() func(*gin.Context) {
	return func(ctx *gin.Context) {
		m := NewGinMiddleware(ctx, &factory.userDataSigningSvc)
		m.SetRole()
	}
}

func (factory *GinMiddlewareFactory) EnforcePolicy(action, object string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		m := NewGinMiddleware(ctx, nil)
		AllowIf(m, func(role string) bool {
			return factory.policies.RoleCan(role, action, object)
		})
	}
}
