package handlers

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"server/db"
	"server/models"

	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DownloadFile(c *fiber.Ctx) error {
	fileIDParam := c.Params("id")
	fileID, err := strconv.ParseUint(fileIDParam, 10, 32)
	if err != nil {
		log.Printf("Invalid file ID format: %s, Error: %v\n", fileIDParam, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID format",
		})
	}

	var dbFile models.File
	result := db.GetDB().First(&dbFile, uint(fileID))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("File with ID %d not found\n", fileID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found",
			})
		}
		log.Printf("Error retrieving file ID %d from database: %v\n", fileID, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve file from database",
		})
	}

	contentType := dbFile.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Set(fiber.HeaderContentType, contentType)

	encodedFilename := url.PathEscape(dbFile.FileName)
	disposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", dbFile.FileName, encodedFilename)
	c.Set(fiber.HeaderContentDisposition, disposition)

	c.Set(fiber.HeaderContentLength, strconv.Itoa(len(dbFile.Data)))

	log.Printf("Serving file ID %d: %s (%s)\n", dbFile.ID, dbFile.FileName, contentType)
	return c.Send(dbFile.Data)
}
