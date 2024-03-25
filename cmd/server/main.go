package main

import (
	"fmt"
	"go-backend-template/internal/impl"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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
