package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	GUID        string `gorm:"type:uuid;default:gen_random_uuid()"`
	UserId      uint   `gorm:"not null"`
	FileName    string
	ContentType string
	Data        []byte `gorm:"type:bytea"`
}
