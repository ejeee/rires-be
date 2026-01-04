package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah instance global database
var DB *gorm.DB

// Connect membuat koneksi ke MySQL database menggunakan GORM
func Connect(dsn string) error {
	var err error

	// Konfigurasi GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Tampilkan SQL query di console
		NowFunc: func() time.Time {
			return time.Now().Local() // Gunakan timezone lokal
		},
	}

	// Buka koneksi ke database
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// GetDB mengembalikan instance database
func GetDB() *gorm.DB {
	return DB
}

// CloseDB menutup koneksi database
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}