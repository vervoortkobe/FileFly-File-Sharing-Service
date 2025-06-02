package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, pepper string) (string, error) {
	passwordWithPepper := pepper + password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithPepper), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hash string, pepper string) bool {
	passwordWithPepper := pepper + password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordWithPepper))
	return err == nil
}
