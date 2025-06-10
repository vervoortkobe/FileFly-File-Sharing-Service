package db

import (
	"log"
	"server/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	registerEnums()

	log.Println("Running database migrations...")
	err = db.AutoMigrate(&models.File{}, &models.User{}, &models.PasswordResetToken{})
	if err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}
	log.Println("Database migration completed successfully.")
}

func registerEnums() {
	log.Println("Creating ENUM types in PostgreSQL...")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	roleEnumExists := false
	err = sqlDB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM pg_type t JOIN pg_namespace n ON (n.oid = t.typnamespace)
		WHERE t.typname = 'user_role' AND n.nspname = current_schema()
	)`).Scan(&roleEnumExists)

	if err != nil {
		log.Fatalf("Failed to check if user_role enum exists: %v", err)
	}
	if !roleEnumExists {
		err = db.Exec(`CREATE TYPE user_role AS ENUM ('guest', 'user', 'admin')`).Error
		if err != nil {
			log.Fatalf("Failed to create user_role enum: %v", err)
		}
		log.Println("Created user_role ENUM type")
	}

	tierEnumExists := false
	err = sqlDB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM pg_type t JOIN pg_namespace n ON (n.oid = t.typnamespace)
		WHERE t.typname = 'user_tier' AND n.nspname = current_schema()
	)`).Scan(&tierEnumExists)

	if err != nil {
		log.Fatalf("Failed to check if user_tier enum exists: %v", err)
	}
	if !tierEnumExists {
		err = db.Exec(`CREATE TYPE user_tier AS ENUM ('guest', 'free', 'premium')`).Error
		if err != nil {
			log.Fatalf("Failed to create user_tier enum: %v", err)
		}
		log.Println("Created user_tier ENUM type")
	}
}
