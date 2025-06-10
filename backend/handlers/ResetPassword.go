package handlers

import (
	"errors"
	"log"
	"server/auth"
	"server/db"
	"server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ResetPasswordPayload struct {
	Token       string `json:"token" form:"token" validate:"required"`
	NewPassword string `json:"newPassword" form:"newPassword" validate:"required,min=8"`
}

func ResetPassword(c *fiber.Ctx) error {
	payload := new(ResetPasswordPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON or form data",
		})
	}

	if payload.Token == "" || payload.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token and new password are required.",
		})
	}
	if len(payload.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters long.",
		})
	}

	var passwordResetToken models.PasswordResetToken
	result := db.GetDB().Where("token = ?", payload.Token).Preload("User").First(&passwordResetToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Invalid or expired password reset token used: %s\n", payload.Token)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or expired password reset token.",
			})
		}
		log.Printf("Error finding password reset token %s: %v\n", payload.Token, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error validating token.",
		})
	}

	if time.Now().After(passwordResetToken.ExpiresAt) {
		log.Printf("Expired password reset token used: %s for user ID %d\n", payload.Token, passwordResetToken.UserID)
		db.GetDB().Delete(&passwordResetToken)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password reset token has expired.",
		})
	}
	hashedPassword, err := auth.HashPassword(payload.NewPassword, db.GetPasswordPepper())
	if err != nil {
		log.Printf("Error hashing new password for user ID %d: %v\n", passwordResetToken.UserID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error processing new password.",
		})
	}

	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", passwordResetToken.UserID).Update("password", hashedPassword).Error; err != nil {
			return err
		}

		if err := tx.Delete(&passwordResetToken).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("Error updating password or deleting token for user ID %d: %v\n", passwordResetToken.UserID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update password.",
		})

	}

	log.Printf("Password reset for user ID %d\n", passwordResetToken.UserID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successful.",
	})
}
