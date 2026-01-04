package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct menyimpan semua konfigurasi aplikasi
type Config struct {
	// Server config
	AppName string
	AppEnv  string
	AppPort string

	// Database config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT config
	JWTSecret       string
	JWTExpiredHours string

	// External API config
	APIBaseURL   string
	APIToken string
}

// Global variable untuk config
var AppConfig *Config

// LoadConfig membaca file .env dan menginisialisasi config
func LoadConfig() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Buat instance config
	AppConfig = &Config{
		AppName: getEnv("APP_NAME", "Golang API"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "student_rires2"),

		JWTSecret:       getEnv("JWT_SECRET", ""),
		JWTExpiredHours: getEnv("JWT_EXPIRED_HOURS", "24"),

		APIBaseURL:   getEnv("API_URL", ""),
		APIToken: getEnv("API_TOKEN", ""),
	}

	return nil
}

// getEnv membaca environment variable dengan fallback ke default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDSN mengembalikan Data Source Name untuk MySQL connection
func (c *Config) GetDSN() string {
	// Format: username:password@tcp(host:port)/dbname?params
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}