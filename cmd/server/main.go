package main

import (
	"fmt"
	"go-backend-template/internal/impl"
	"go-backend-template/internal/util"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	v1 := r.Group("/v1")

	v1.GET("/posts/:post", util.DecorateHandler(func(ctx *gin.Context) (any, error) {
		post, ok := util.GetParamInt(&ctx.Params, "post")

		if !ok {
			return nil, &util.RequestError{
				StatusCode: 500,
				Err:        fmt.Errorf("post not specified"),
			}
		}

		if post == 5 {
			type response struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			}

			return response{
				Title:   "Foobar",
				Content: "Baaz",
			}, nil
		}

		return nil, &util.RequestError{
			StatusCode: 404,
			Err:        fmt.Errorf("post not found"),
		}
	}))

	r.Run("0.0.0.0:8081")

	pg, err := impl.NewPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_PORT"),
	)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = pg.Migrate("./assets/migrations")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	userRepo := impl.NewPGUserRepo(&pg)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	authSvc := impl.NewDefaultAuthSvc(userRepo, "foobar")
	// err = authSvc.CreateUser("admin@example.com", "lorem ipsum", "admin")

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	user, err := authSvc.AuthenticateUser("admin@example.com", "lorem ipsum")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Successfully authenticated", user)

	signingSvc := impl.NewJWTSigningSvc("foo")
	token, err := signingSvc.Sign(user.ID, user.Role)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Successfully signed", token)

	parsedID, parsedRole, err := signingSvc.Parse(token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Successfully parsed", parsedID, parsedRole)
}
