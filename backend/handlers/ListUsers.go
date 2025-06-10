package handlers

import (
	"log"
	"server/db"
	"server/models"

	"github.com/gofiber/fiber/v2"
)

func ListUsers(c *fiber.Ctx) error {
	var users []models.User

	result := db.GetDB().Model(&models.User{}).
		Order("created_at desc").
		Select("id, email, role, tier, created_at").
		Find(&users)

	if result.Error != nil {
		log.Println("Error listing users from database:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not list users",
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusOK).JSON([]models.User{})
	}

	log.Printf("Listed %d users\n", len(users))
	return c.Status(fiber.StatusOK).JSON(users)
}
