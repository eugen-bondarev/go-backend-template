package main

import (
	"fmt"
	"go-backend-template/internal/impl"
	"go-backend-template/internal/middleware"
	"go-backend-template/internal/model"
	"go-backend-template/internal/util"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

type Permission struct {
	Object  string
	Action  string
	Subject string
}

func NewPermission(subject, action, object string) Permission {
	return Permission{
		Subject: subject,
		Action:  action,
		Object:  object,
	}
}

type Permissions struct {
	permissions []Permission
}

func NewPermissions() Permissions {
	return Permissions{
		permissions: make([]Permission, 0),
	}
}

func (p *Permissions) Add(subject, action, object string) {
	p.permissions = append(p.permissions, NewPermission(subject, action, object))
}

func (p *Permissions) RoleCan(subject, action, object string) bool {
	for _, permission := range p.permissions {
		if permission.Subject == subject && permission.Action == action && permission.Object == object {
			return true
		}
	}
	return false
}

type App struct {
	userRepo    model.UserRepo
	signingSvc  model.SigningSvc
	authSvc     model.AuthSvc
	permissions Permissions
}

func NewApp() (App, error) {
	pg, err := impl.NewPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_PORT"),
	)

	if err != nil {
		return App{}, err
	}

	err = pg.Migrate("./assets/migrations")

	if err != nil {
		return App{}, err
	}

	userRepo := impl.NewPGUserRepo(&pg)
	authSvc := impl.NewDefaultAuthSvc(userRepo, "foobar")
	signingSvc := impl.NewJWTSigningSvc("foo")

	permissions := NewPermissions()
	permissions.Add("admin", "list", "users")
	permissions.Add("user", "list", "users")

	return App{
		userRepo:    userRepo,
		signingSvc:  signingSvc,
		authSvc:     authSvc,
		permissions: permissions,
	}, nil
}

func (app *App) users() ([]model.User, error) {
	return app.userRepo.GetUsers()
}

func (app *App) authMiddleware(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	middleware.Auth(
		app.signingSvc,
		authHeader,
		func(ID int, role string) {
			ctx.Set("ID", ID)
			ctx.Set("role", role)
		},
	)
}

func (app *App) requiredAuthMiddleware(action, object string) func(*gin.Context) error {
	return func(ctx *gin.Context) error {
		role := ctx.GetString("role")

		if !app.permissions.RoleCan(role, action, object) {
			return &util.RequestError{
				StatusCode: 403,
				Err:        fmt.Errorf("unauthorized"),
			}
		}

		return nil
	}
}

func (app *App) decorateWithPermissions(action, object string, handler func(*gin.Context) (any, error)) gin.HandlersChain {
	chain := make(gin.HandlersChain, 0, 3)
	chain = append(chain, util.DecorateMiddleware(app.authMiddleware))
	chain = append(chain, util.DecorateRequiredMiddleware(app.requiredAuthMiddleware(action, object)))
	chain = append(chain, util.DecorateHandler(handler))
	return chain
}

func (app *App) decorateForAnyone(handler func(*gin.Context) (any, error)) gin.HandlersChain {
	chain := make(gin.HandlersChain, 0, 1)
	chain = append(chain, util.DecorateHandler(handler))
	return chain
}

func main() {
	godotenv.Load()

	app, err := NewApp()

	if err != nil {
		panic(err.Error())
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	v1 := r.Group("/v1")

	v1.GET(
		"/users",
		app.decorateWithPermissions(
			"list",
			"users",
			func(ctx *gin.Context) (any, error) {
				return app.users()
			},
		)...,
	)

	auth := v1.Group("/auth")

	auth.POST(
		"/login",
		app.decorateForAnyone(
			func(ctx *gin.Context) (any, error) {
				type Payload struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}
				var payload Payload
				ctx.ShouldBindBodyWith(&payload, binding.JSON)

				return app.login(payload.Email, payload.Password)
			},
		)...,
	)

	auth.POST(
		"/register",
		app.decorateForAnyone(
			func(ctx *gin.Context) (any, error) {
				type Payload struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}
				var payload Payload
				ctx.ShouldBindBodyWith(&payload, binding.JSON)

				return nil, app.register(payload.Email, payload.Password)
			},
		)...,
	)

	r.Run("0.0.0.0:4200")
}
