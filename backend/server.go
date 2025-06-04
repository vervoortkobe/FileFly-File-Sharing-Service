package main

import (
	"fmt"
	"log"
	"os"

	"server/auth"
	"server/db"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables where possible.")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	db.Connect()
	defer db.Close()
	db.Migrate(db.GetDB())

	app := fiber.New()

	app.Post("/register", auth.HandleRegisterUser)
	app.Post("/login", auth.HandleLoginUser)

	fmt.Printf("[⚡] WebServer listening on [http://localhost:%s]!\n", PORT)
	log.Fatal(app.Listen(":" + PORT))
}
