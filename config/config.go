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

	// NEOMAA Database config (data mahasiswa)
	DBNeomaaHost     string
	DBNeomaaPort     string
	DBNeomaaUser     string
	DBNeomaaPassword string
	DBNeomaaName     string

	// NEOMAAREF Database config (data referensi: prodi, fakultas)
	DBNeomaaRefHost     string
	DBNeomaaRefPort     string
	DBNeomaaRefUser     string
	DBNeomaaRefPassword string
	DBNeomaaRefName     string

	// SIMPEG Database config (data pegawai)
	DBSimpegHost     string
	DBSimpegPort     string
	DBSimpegUser     string
	DBSimpegPassword string
	DBSimpegName     string

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

		// NEOMAA Database
		DBNeomaaHost:     getEnv("DB_NEOMAA_HOST", "localhost"),
		DBNeomaaPort:     getEnv("DB_NEOMAA_PORT", "3306"),
		DBNeomaaUser:     getEnv("DB_NEOMAA_USER", "root"),
		DBNeomaaPassword: getEnv("DB_NEOMAA_PASSWORD", ""),
		DBNeomaaName:     getEnv("DB_NEOMAA_NAME", "neomaa"),

		// NEOMAAREF Database
		DBNeomaaRefHost:     getEnv("DB_NEOMAAREF_HOST", "localhost"),
		DBNeomaaRefPort:     getEnv("DB_NEOMAAREF_PORT", "3306"),
		DBNeomaaRefUser:     getEnv("DB_NEOMAAREF_USER", "root"),
		DBNeomaaRefPassword: getEnv("DB_NEOMAAREF_PASSWORD", ""),
		DBNeomaaRefName:     getEnv("DB_NEOMAAREF_NAME", "neomaaref"),

		// SIMPEG Database
		DBSimpegHost:     getEnv("DB_SIMPEG_HOST", "localhost"),
		DBSimpegPort:     getEnv("DB_SIMPEG_PORT", "3306"),
		DBSimpegUser:     getEnv("DB_SIMPEG_USER", "root"),
		DBSimpegPassword: getEnv("DB_SIMPEG_PASSWORD", ""),
		DBSimpegName:     getEnv("DB_SIMPEG_NAME", "newsimpeg"),

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

// GetDSNNeomaa mengembalikan DSN untuk database NEOMAA
func (c *Config) GetDSNNeomaa() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBNeomaaUser,
		c.DBNeomaaPassword,
		c.DBNeomaaHost,
		c.DBNeomaaPort,
		c.DBNeomaaName,
	)
}

// GetDSNNeomaaRef mengembalikan DSN untuk database NEOMAAREF
func (c *Config) GetDSNNeomaaRef() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBNeomaaRefUser,
		c.DBNeomaaRefPassword,
		c.DBNeomaaRefHost,
		c.DBNeomaaRefPort,
		c.DBNeomaaRefName,
	)
}

// GetDSNSimpeg mengembalikan DSN untuk database SIMPEG
func (c *Config) GetDSNSimpeg() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBSimpegUser,
		c.DBSimpegPassword,
		c.DBSimpegHost,
		c.DBSimpegPort,
		c.DBSimpegName,
	)
}