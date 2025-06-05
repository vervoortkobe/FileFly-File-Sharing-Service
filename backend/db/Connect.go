package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db             *gorm.DB
	passwordPepper string
	err            error
)

func Connect() {
	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the database.")
}

func Close() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Printf("Error closing database connection: %v", err)
		return
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("Database connection is not initialized. Call Connect() first.")
	}
	return db
}

func GetPasswordPepper() string {
	passwordPepper = os.Getenv("PASSWORD_PEPPER")
	if passwordPepper == "" {
		log.Fatal("FATAL: PASSWORD_PEPPER environment variable not set. Application cannot start.")
	}
	return passwordPepper
}
