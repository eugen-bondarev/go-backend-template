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

type App struct {
	userRepo    model.UserRepo
	signingSvc  model.SigningSvc
	authSvc     model.AuthSvc
	permissions []Permission
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

	permissions := make([]Permission, 0)
	permissions = append(permissions, NewPermission("user", "list", "users"))
	permissions = append(permissions, NewPermission("admin", "list", "users"))

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
		role, _ := ctx.Get("role")

		for _, permission := range app.permissions {
			if permission.Subject == role && permission.Action == action && permission.Object == object {
				return nil
			}
		}

		return &util.RequestError{
			StatusCode: 403,
			Err:        fmt.Errorf("unauthorized"),
		}
	}
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
		util.DecorateMiddleware(app.authMiddleware),
		util.DecorateRequiredMiddleware(app.requiredAuthMiddleware("list", "users")),
		util.DecorateHandler(func(ctx *gin.Context) (any, error) {
			return app.users()
		}),
	)

	auth := v1.Group("/auth")

	auth.POST("/login", util.DecorateHandler(func(ctx *gin.Context) (any, error) {
		type Payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var payload Payload
		ctx.ShouldBindBodyWith(&payload, binding.JSON)

		return app.login(payload.Email, payload.Password)
	}))

	auth.POST("/register", util.DecorateHandler(func(ctx *gin.Context) (any, error) {
		type Payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var payload Payload
		ctx.ShouldBindBodyWith(&payload, binding.JSON)

		return nil, app.register(payload.Email, payload.Password)
	}))

	r.Run("0.0.0.0:8081")
}
