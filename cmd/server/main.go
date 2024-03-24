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

	err = pg.Migrate()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	userRepo := impl.NewPGUserRepo(&pg)

	fmt.Println(userRepo)
}
