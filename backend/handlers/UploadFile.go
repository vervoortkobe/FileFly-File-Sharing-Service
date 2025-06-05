package handlers

import (
	"io"
	"log"

	"server/db"
	"server/models"

	"github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file from form:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to retrieve file from form",
		})
	}

	uploadedFile, err := fileHeader.Open()
	if err != nil {
		log.Println("Error opening uploaded file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer uploadedFile.Close()

	fileBytes, err := io.ReadAll(uploadedFile)
	if err != nil {
		log.Println("Error reading file content:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file content",
		})
	}

	dbFile := models.File{
		UserId: uint(0),
		//UserId:    c.Locals("user_id").(uint),
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Data:        fileBytes,
	}

	result := db.GetDB().Create(&dbFile)
	if result.Error != nil {
		log.Println("Error saving file to database:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file to database",
		})
	}

	log.Printf("File '%s' uploaded successfully. ID: %d\n", dbFile.FileName, dbFile.ID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "File uploaded successfully",
		"id":       dbFile.ID,
		"filename": dbFile.FileName,
	})
}
