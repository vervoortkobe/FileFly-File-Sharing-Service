package handlers

import (
	"log"
	"strconv"

	"server/db"
	"server/models"

	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserByID(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		log.Printf("Invalid user ID format: %s, Error: %v\n", userIDParam, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var dbUser models.User
	result := db.GetDB().First(&dbUser, uint(userID))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with ID %d not found\n", userID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		log.Printf("Error retrieving user ID %d from database: %v\n", userID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve user from database",
		})
	}

	log.Printf("Listed user with ID %d: %s\n", dbUser.ID, dbUser.Email)
	return c.Status(fiber.StatusOK).JSON(dbUser)
}
