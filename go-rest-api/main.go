// main.go

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	a := App{}
	a.Initialize(
		"localhost",
		os.Getenv("PORT"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("SSL_MODE"))

	a.Run(":8080")
}
