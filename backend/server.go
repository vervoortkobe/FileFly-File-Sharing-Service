package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"server/auth"
	"server/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db             *gorm.DB
	passwordPepper string
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

	passwordPepper = os.Getenv("PASSWORD_PEPPER")
	if passwordPepper == "" {
		log.Fatal("FATAL: PASSWORD_PEPPER environment variable not set. Application cannot start.")
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		dbHost,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	log.Println("Creating ENUM types in PostgreSQL...")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	roleEnumExists := false
	err = sqlDB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM pg_type t JOIN pg_namespace n ON (n.oid = t.typnamespace)
		WHERE t.typname = 'user_role' AND n.nspname = current_schema()
	)`).Scan(&roleEnumExists)

	if err != nil {
		log.Fatalf("Failed to check if user_role enum exists: %v", err)
	}
	if !roleEnumExists {
		err = db.Exec(`CREATE TYPE user_role AS ENUM ('guest', 'user', 'admin')`).Error
		if err != nil {
			log.Fatalf("Failed to create user_role enum: %v", err)
		}
		log.Println("Created user_role ENUM type")
	}

	tierEnumExists := false
	err = sqlDB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM pg_type t JOIN pg_namespace n ON (n.oid = t.typnamespace)
		WHERE t.typname = 'user_tier' AND n.nspname = current_schema()
	)`).Scan(&tierEnumExists)

	if err != nil {
		log.Fatalf("Failed to check if user_tier enum exists: %v", err)
	}
	if !tierEnumExists {
		err = db.Exec(`CREATE TYPE user_tier AS ENUM ('guest', 'free', 'premium')`).Error
		if err != nil {
			log.Fatalf("Failed to create user_tier enum: %v", err)
		}
		log.Println("Created user_tier ENUM type")
	}

	log.Println("Running database migrations...")
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully.")

	app := fiber.New()

	app.Post("/register", handleRegisterUser)

	fmt.Printf("[⚡] WebServer listening on [http://localhost:%s]!\n", PORT)
	log.Fatal(app.Listen(":" + PORT))
}

func handleRegisterUser(c *fiber.Ctx) error {
	payload := new(domain.RegisterUser)

	if err := c.BodyParser(payload); err != nil {
		log.Printf("Error parsing /register request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if payload.Email == "" || payload.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}
	if len(payload.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
	}
	if !strings.Contains(payload.Email, "@") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	hashedPassword, err := auth.HashPassword(payload.Password, passwordPepper)
	if err != nil {
		log.Printf("Error hashing password for %s: %v\n", payload.Email, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing registration"})
	}

	user := domain.User{
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     domain.GuestRole,
		Tier:     domain.GuestTier,
	}

	result := db.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating user %s in DB: %v\n", payload.Email, result.Error)
		if strings.Contains(strings.ToLower(result.Error.Error()), "unique constraint") ||
			strings.Contains(strings.ToLower(result.Error.Error()), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
	}

	responseUser := fiber.Map{
		"id":         user.ID,
		"guid":       user.GUID,
		"email":      user.Email,
		"role":       user.Role.String(),
		"tier":       user.Tier.String(),
		"created_at": user.CreatedAt,
	}
	return c.Status(fiber.StatusCreated).JSON(responseUser)
}

func handleLoginUser(c *fiber.Ctx) error {
	payload := new(domain.LoginUser)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user domain.User
	result := db.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
		}
		log.Printf("Error fetching user %s for login: %v\n", payload.Email, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Login failed"})
	}

	if !auth.CheckPasswordHash(payload.Password, user.Password, passwordPepper) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user_id": user.ID,
		"guid":    user.GUID,
		"email":   user.Email,
		"role":    user.Role.String(),
	})
}
