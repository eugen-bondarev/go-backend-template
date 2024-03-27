package main

import (
	"go-backend-template/internal/impl"
	"go-backend-template/internal/model"
	"go-backend-template/internal/util"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
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
	authSvc := impl.NewDefaultAuthSvc(userRepo, "foobar")
	signingSvc := impl.NewJWTSigningSvc("foo")

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

func main() {
	godotenv.Load()

	app, err := NewApp()

	if err != nil {
		panic(err.Error())
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (any, error) {
				return "world", nil
			},
		},
		"users": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "User",
				Fields: graphql.Fields{
					"id": &graphql.Field{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "User ID",
					},
					"email": &graphql.Field{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "User email",
					},
					"role": &graphql.Field{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "User role",
					},
				},
			})),
			Args: graphql.FieldConfigArgument{
				"role": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				var users []model.User
				var err error

				if p.Args["role"] != nil {
					users, err = app.userRepo.GetUsersByRole(p.Args["role"].(string))
				} else {
					users, err = app.userRepo.GetUsers()
				}

				if err != nil {
					return nil, nil
				}

				output := util.Map(users, func(user model.User) model.APIUser {
					return model.APIUser{ID: user.ID, Email: user.Email, Role: user.Role}
				})

				return output, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	httpRouter := gin.Default()
	httpRouter.GET("/apollo", func(ctx *gin.Context) {
		ctx.File("./assets/html/apollo-graphql.html")
	})
	httpRouter.POST("/graphql", func(ctx *gin.Context) {
		type payload struct {
			Query string `json:"query"`
		}
		var pl payload
		ctx.ShouldBindBodyWith(&pl, binding.JSON)

		params := graphql.Params{Schema: schema, RequestString: pl.Query}
		r := graphql.Do(params)
		if len(r.Errors) > 0 {
			log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
		}
		// rJSON, _ := json.Marshal(r)
		ctx.JSON(200, r)
		// fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
		// return
	})

	httpRouter.Run("0.0.0.0:4200")

	// mw := middleware.NewGinMiddlewareFactory(
	// 	app.signingSvc,
	// 	&app.policies,
	// )

	// r := gin.Default()
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
	// 		type Payload struct {
	// 			Email    string `json:"email"`
	// 			Password string `json:"password"`
	// 		}
	// 		var payload Payload
	// 		ctx.ShouldBindBodyWith(&payload, binding.JSON)

	// 		return app.login(payload.Email, payload.Password)
	// 	}),
	// )

	// auth.POST(
	// 	"/register",
	// 	util.DecorateHandler(func(ctx *gin.Context) (any, error) {
	// 		type Payload struct {
	// 			Email    string `json:"email"`
	// 			Password string `json:"password"`
	// 		}
	// 		var payload Payload
	// 		ctx.ShouldBindBodyWith(&payload, binding.JSON)

	// 		return nil, app.register(payload.Email, payload.Password)
	// 	}),
	// )

	// r.Run("0.0.0.0:4200")
}
