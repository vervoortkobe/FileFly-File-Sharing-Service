package models

import (
	"time"

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

type FileInfo struct {
	ID          uint      `json:"id"`
	FileName    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	CreatedAt   time.Time `json:"createdAt"`
}
