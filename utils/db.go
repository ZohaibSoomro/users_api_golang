package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zohaibsoomro/users_api_golang/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDb() *gorm.DB {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Build DSN using environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.Config.Logger = logger.Default.LogMode(logger.Silent)

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate schema: " + err.Error())
	}

	return db
}
