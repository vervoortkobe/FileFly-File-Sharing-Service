package handlers

import (
	"fmt"
	"log"
	"time"

	"server/db"
	"server/mail"
	"server/models"

	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestPasswordResetPayload struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

func RequestPasswordReset(c *fiber.Ctx) error {
	payload := new(RequestPasswordResetPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON or form data",
		})
	}

	if payload.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	var user models.User
	result := db.GetDB().Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Password reset requested for non-existent email: %s\n", payload.Email)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "If an account with that email exists, a password reset link has been sent.",
			})
		}
		log.Printf("Error finding user by email %s: %v\n", payload.Email, result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error processing your request.",
		})
	}

	resetToken := uuid.NewString()
	tokenExpiration := time.Now().Add(1 * time.Hour)

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.PasswordResetToken{}).Error; err != nil {
			return err
		}

		newPasswordResetToken := models.PasswordResetToken{
			UserID:    user.ID,
			Token:     resetToken,
			ExpiresAt: tokenExpiration,
		}
		if err := tx.Create(&newPasswordResetToken).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("Error saving password reset token for user %d: %v\n", user.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not save password reset token.",
		})
	}

	resetURL := fmt.Sprintf("%s/reset-password.html?token=%s", c.BaseURL(), resetToken)

	emailSubject := "Password Reset Request"
	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>You recently requested to reset your password for your account.</p>
		<p>Click the link below to reset it:</p>
		<p><a href="%s">%s</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you did not request a password reset, please ignore this email or contact support if you have concerns.</p>
		<p>Thanks,</p>
		<p>Your Application Team</p>
	`, user.Email, resetURL, resetURL)

	go func(recipientName, recipientAddress, subject, body, token string) {
		err := mail.SendEmail(mail.MailConfig{
			SMTPServer:   "smtp.example.com",
			SMTPPort:     587,
			SMTPUser:     "user@example.com",
			SMTPPassword: "yourpassword",
			FromEmail:    "no-reply@yourapp.com",
			FromName:     "Your Application",
		}, recipientName, recipientAddress, subject, body)

		if err != nil {
			log.Printf("Failed to send password reset email to %s (token: %s): %v\n", recipientAddress, token, err)
		} else {
			log.Printf("Password reset email sent to %s (token: %s)\n", recipientAddress, token)
		}
	}(user.Email, user.Email, emailSubject, emailBody, resetToken)

	log.Printf("Password reset process initiated for user: %s (ID: %d)\n", user.Email, user.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If an account with that email exists, a password reset link has been sent.",
	})
}
