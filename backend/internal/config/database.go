package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
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
		log.Printf("Error connecting to database: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		log.Printf("Error pinging database: %v", err)
		return nil
	}

	log.Println("Successfully connected to the database")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
		&domain.Plant{},
		&domain.Disease{},
		&domain.Inspection{},
	)

	if err != nil {
		log.Printf("Error during auto migration: %v", err)
		return nil
	}
	log.Printf("Database auto migration completed successfully")
	return db
}
