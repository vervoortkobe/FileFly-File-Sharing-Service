package auth

import (
	"errors"
	"log"
	"server/db"
	"server/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HandleLoginUser(c *fiber.Ctx) error {
	payload := new(domain.LoginUser)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user domain.User
	result := db.GetDB().Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
		}
		log.Printf("Error fetching user %s for login: %v\n", payload.Email, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Login failed"})
	}

	if !CheckPasswordHash(payload.Password, user.Password, db.GetPasswordPepper()) {
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
