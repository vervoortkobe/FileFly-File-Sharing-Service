package handlers

import (
	"log"
	"server/db"
	"server/models"

	"github.com/gofiber/fiber/v2"
)

func ListFiles(c *fiber.Ctx) error {
	var files []models.File

	result := db.GetDB().Model(&models.File{}).
		Order("created_at desc").
		Select("id, file_name, content_type, created_at").
		Find(&files)

	if result.Error != nil {
		log.Println("Error listing files from database:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list files",
		})
	}

	if len(files) == 0 {
		return c.Status(fiber.StatusOK).JSON([]models.File{})
	}

	log.Printf("Listed %d files\n", len(files))
	return c.Status(fiber.StatusOK).JSON(files)
}
