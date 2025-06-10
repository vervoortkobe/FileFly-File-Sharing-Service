package main

import (
	"fmt"
	"log"
	"os"

	"server/auth"
	"server/db"
	"server/handlers"
	"server/models"

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
		var fileCount int64
		db.GetDB().Model(&models.File{}).Count(&fileCount)
		return c.SendString("Files: " + fmt.Sprintf("%d", fileCount))
	})

	api := app.Group("/api")

	api.Post("/register", auth.HandleRegisterUser)
	api.Post("/login", auth.HandleLoginUser)

	api.Get("/users", handlers.ListUsers)
	api.Get("/files", handlers.ListFiles)
	api.Post("/upload", handlers.UploadFile)
	api.Get("/download/:id", handlers.DownloadFile)

	fmt.Printf("[⚡] WebServer listening on [http://localhost:%s]!\n", PORT)
	log.Fatal(app.Listen(":" + PORT))
}
