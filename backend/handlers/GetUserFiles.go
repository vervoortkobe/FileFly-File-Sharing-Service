package handlers

import (
	"errors"
	"log"
	"server/db"
	"server/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserFiles(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		log.Printf("Invalid UserID format for GetUserFiles: %s, Error: %v\n", userIDParam, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UserID format",
		})
	}

	var user models.User
	if err := db.GetDB().First(&user, uint(userID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User with ID %d not found when trying to list their files\n", userID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		log.Printf("Error checking for user ID %d: %v\n", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error verifying user"})
	}

	var userFiles []models.File
	result := db.GetDB().Model(&models.File{}).
		Where("user_id = ?", uint(userID)).
		Order("created_at desc").
		Select("id, file_name, content_type, created_at, user_id, guid").
		Find(&userFiles)

	if result.Error != nil {
		log.Printf("Error listing files for user ID %d: %v\n", userID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list files for the user",
		})
	}

	if len(userFiles) == 0 {
		log.Printf("No files found for user ID %d\n", userID)
		return c.Status(fiber.StatusOK).JSON([]models.File{})
	}

	log.Printf("Listed %d files for user ID %d\n", len(userFiles), userID)
	return c.Status(fiber.StatusOK).JSON(userFiles)
}
