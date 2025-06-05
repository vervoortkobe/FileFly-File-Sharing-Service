package main

import (
	"fmt"
	"log"
	"os"

	"server/auth"
	"server/db"
	"server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})

	app.Use(logger.New())

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	api := app.Group("/api")

	api.Post("/register", auth.HandleRegisterUser)
	api.Post("/login", auth.HandleLoginUser)

	api.Post("/upload", handlers.UploadFile)
	api.Get("/download/:id", handlers.DownloadFile)
	api.Get("/files", handlers.ListFiles)

	fmt.Printf("[⚡] WebServer listening on [http://localhost:%s]!\n", PORT)
	log.Fatal(app.Listen(":" + PORT))
}
