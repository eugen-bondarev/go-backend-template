package main

import (
	"go-backend-template/internal/impl"
	"go-backend-template/internal/model"
	"log"
	"net/http"
	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type App struct {
	userRepo   model.UserRepo
	signingSvc model.SigningSvc
	authSvc    model.AuthSvc
	policies   impl.Policies
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
	authSvc := impl.NewDefaultAuthSvc(userRepo, os.Getenv("PEPPER"))
	signingSvc := impl.NewJWTSigningSvc(os.Getenv("JWT_SECRET"))

	policies := impl.NewPolicies()
	policies.Add("admin", "index", "users")
	policies.Add("admin", "manage", "users")

	return App{
		userRepo:   userRepo,
		signingSvc: signingSvc,
		authSvc:    authSvc,
		policies:   policies,
	}, nil
}

type query struct{}

func (query) Hello() string {
	return "hello, wordl"
}

func main() {
	s := `
		type Query {
			hello: String!
		}
	`
	schema := graphql.MustParseSchema(s, &query{})
	// r := gin.Default()

	http.Handle("/", &relay.Handler{Schema: schema})
	log.Fatal(http.ListenAndServe(":4200", nil))
	// r.GET("/query", &relay.Handler{Schema: schema})

	// godotenv.Load()

	// app, err := NewApp()

	// if err != nil {
	// 	panic(err.Error())
	// }

	// if os.Getenv("GIN_MODE") == "release" {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	// mw := middleware.NewGinMiddlewareFactory(
	// 	app.signingSvc,
	// 	&app.policies,
	// )

	// r := gin.Default()
	// r.SetTrustedProxies(nil)
	// r.Use(util.GinConfigureCors(os.Getenv("CORS_ALLOW_ORIGINS")))

	// util.GinHealthz(r)

	// v1 := r.Group("/v1")

	// v1.GET(
	// 	"/users",
	// 	mw.SetRole(),
	// 	mw.EnforcePolicy("index", "users"),
	// 	util.DecorateHandler(func(ctx *gin.Context) (any, error) {
	// 		return app.userRepo.GetUsers()
	// 	}),
	// )

	// v1.DELETE(
	// 	"/users/:id",
	// 	mw.SetRole(),
	// 	mw.EnforcePolicy("manage", "users"),
	// 	util.DecorateHandler(func(ctx *gin.Context) (any, error) {
	// 		id, err := strconv.Atoi(ctx.Params.ByName("id"))

	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		return nil, app.userRepo.DeleteUserByID(id)
	// 	}),
	// )

	// auth := v1.Group("/auth")

	// auth.POST(
	// 	"/login",
	// 	util.DecorateHandler(func(ctx *gin.Context) (any, error) {
	// 		payload, err := util.GinGetBody[struct {
	// 			dto.WithEmail
	// 			dto.WithPassword
	// 		}](ctx)

	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		return app.login(payload.Email, payload.Password)
	// 	}),
	// )

	// auth.POST(
	// 	"/register",
	// 	util.DecorateHandler(func(ctx *gin.Context) (any, error) {
	// 		payload, err := util.GinGetBody[struct {
	// 			dto.WithEmail
	// 			dto.WithPassword
	// 		}](ctx)

	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		return nil, app.register(payload.Email, payload.Password)
	// 	}),
	// )

	// r.Run("0.0.0.0:4200")
}
