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
		fmt.Sprintf(
			"user=%s host=%s port=%s password=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_PASS"),
		),
	)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(pg)
}
