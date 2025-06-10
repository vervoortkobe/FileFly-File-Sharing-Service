package models

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	GUID        string `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID      uint   `gorm:"not null"`
	FileName    string
	ContentType string
	Data        []byte `gorm:"type:bytea"`
}
