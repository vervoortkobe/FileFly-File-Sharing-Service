package auth

import (
	"log"
	"server/db"
	"server/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandleRegisterUser(c *fiber.Ctx) error {
	payload := new(models.RegisterUser)

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

	hashedPassword, err := HashPassword(payload.Password, db.GetPasswordPepper())
	if err != nil {
		log.Printf("Error hashing password for %s: %v\n", payload.Email, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing registration"})
	}

	user := models.User{
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     models.GuestRole,
		Tier:     models.GuestTier,
	}

	result := db.GetDB().Create(&user)
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
