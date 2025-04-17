package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"user-service/internal/model"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate các bảng
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	err = DB.AutoMigrate(&model.EmailVerificationToken{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}
}

func GetJWTSecret() []byte {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY not set in environment")
	}
	return []byte(secret)
}

func GetString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: %s is not set in environment", key)
	}
	return value
}
