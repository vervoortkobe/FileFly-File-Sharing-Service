package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"server/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	app := fiber.New()

	dsn := "host=localhost user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	app.Post("/register", func(c *fiber.Ctx) error {
		payload := domain.RegisterUser{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		user := domain.User{ID: 1, Email: payload.Email, Password: payload.Password, Role: domain.GuestRole, Tier: domain.GuestTier, CreatedAt: time.Now()}

		result := db.Create(&user) // pass pointer of data to Create
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating user: " + result.Error.Error())
		}
		return c.SendString(fmt.Sprintf("Registered user %s<br>ID: %d", user.Email, user.ID))
	})

	fmt.Printf("[⚡] WebServer listening on [http://localhost:%s]!\n", PORT)
	log.Fatal(app.Listen(":" + PORT))
}
