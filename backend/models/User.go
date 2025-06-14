package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GUID     string `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex" json:"guid"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     Role   `gorm:"type:user_role;not null;default:guest" json:"role"`
	Tier     Tier   `gorm:"type:user_tier;not null;default:guest" json:"tier"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
