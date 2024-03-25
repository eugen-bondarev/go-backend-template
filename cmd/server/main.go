package main

import (
	"fmt"
	"go-backend-template/internal/impl"
	"go-backend-template/internal/model"
	"go-backend-template/internal/util"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

type App struct {
	userRepo   model.UserRepo
	signingSvc model.SigningSvc
	authSvc    model.AuthSvc
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

	if err != nil {
		return App{}, err
	}

	authSvc := impl.NewDefaultAuthSvc(userRepo, "foobar")
	// err = authSvc.CreateUser("admin@example.com", "lorem ipsum", "admin")

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	user, err := authSvc.AuthenticateUser("admin@example.com", "lorem ipsum")

	if err != nil {
		return App{}, err
	}

	fmt.Println("Successfully authenticated", user)

	signingSvc := impl.NewJWTSigningSvc("foo")
	token, err := signingSvc.Sign(user.ID, user.Role)

	if err != nil {
		return App{}, err
	}

	fmt.Println("Successfully signed", token)

	parsedID, parsedRole, err := signingSvc.Parse(token)

	if err != nil {
		return App{}, err
	}

	fmt.Println("Successfully parsed", parsedID, parsedRole)

	return App{
		userRepo:   userRepo,
		signingSvc: signingSvc,
		authSvc:    authSvc,
	}, nil
}

func (app *App) users() ([]model.User, error) {
	return app.userRepo.GetUsers()
}

func (app *App) login(email, plainTextPassword string) (string, error) {
	user, err := app.authSvc.AuthenticateUser(email, plainTextPassword)

	if err != nil {
		return "", err
	}

	token, err := app.signingSvc.Sign(user.ID, user.Role)

	return token, err
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
	v1.GET("/users", util.DecorateHandler(func(ctx *gin.Context) (any, error) {
		return app.users()
	}))
	v1.POST("/login", util.DecorateHandler(func(ctx *gin.Context) (any, error) {
		type Payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var payload Payload
		ctx.ShouldBindBodyWith(&payload, binding.JSON)

		return app.login(payload.Email, payload.Password)
	}))

	r.Run("0.0.0.0:8081")
}
