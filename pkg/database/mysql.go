package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB instances untuk masing-masing database
var (
	DB           *gorm.DB // Database utama (student_rires2)
	DBNeomaa     *gorm.DB // Database NEOMAA (data mahasiswa)
	DBNeomaaRef  *gorm.DB // Database NEOMAAREF (data referensi: prodi, fakultas)
	DBSimpeg     *gorm.DB // Database NEWSIMPEG (data pegawai)
)

// Connect membuat koneksi ke MySQL database menggunakan GORM
func Connect(dsn string) error {
	var err error

	// Konfigurasi GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent mode - no SQL logs
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// Buka koneksi ke database utama
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect to main database: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Main database connected")

	return nil
}

// ConnectExternal membuat koneksi ke database eksternal (NEOMAA, NEOMAAREF, SIMPEG)
func ConnectExternal(dsnNeomaa, dsnNeomaaRef, dsnSimpeg string) error {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent mode
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// Connect to NEOMAA (data mahasiswa)
	var err error
	DBNeomaa, err = gorm.Open(mysql.Open(dsnNeomaa), config)
	if err != nil {
		return fmt.Errorf("failed to connect to NEOMAA database: %w", err)
	}
	
	sqlDBNeomaa, _ := DBNeomaa.DB()
	sqlDBNeomaa.SetMaxIdleConns(5)
	sqlDBNeomaa.SetMaxOpenConns(20)
	sqlDBNeomaa.SetConnMaxLifetime(time.Hour)
	
	log.Println("✅ NEOMAA database connected")

	// Connect to NEOMAAREF (data referensi)
	DBNeomaaRef, err = gorm.Open(mysql.Open(dsnNeomaaRef), config)
	if err != nil {
		return fmt.Errorf("failed to connect to NEOMAAREF database: %w", err)
	}
	
	sqlDBNeomaaRef, _ := DBNeomaaRef.DB()
	sqlDBNeomaaRef.SetMaxIdleConns(5)
	sqlDBNeomaaRef.SetMaxOpenConns(20)
	sqlDBNeomaaRef.SetConnMaxLifetime(time.Hour)
	
	log.Println("✅ NEOMAAREF database connected")

	// Connect to SIMPEG (data pegawai)
	DBSimpeg, err = gorm.Open(mysql.Open(dsnSimpeg), config)
	if err != nil {
		return fmt.Errorf("failed to connect to SIMPEG database: %w", err)
	}
	
	sqlDBSimpeg, _ := DBSimpeg.DB()
	sqlDBSimpeg.SetMaxIdleConns(5)
	sqlDBSimpeg.SetMaxOpenConns(20)
	sqlDBSimpeg.SetConnMaxLifetime(time.Hour)
	
	log.Println("✅ SIMPEG database connected")

	return nil
}

// GetDB mengembalikan instance database utama
func GetDB() *gorm.DB {
	return DB
}

// GetDBNeomaa mengembalikan instance database NEOMAA
func GetDBNeomaa() *gorm.DB {
	return DBNeomaa
}

// GetDBNeomaaRef mengembalikan instance database NEOMAAREF
func GetDBNeomaaRef() *gorm.DB {
	return DBNeomaaRef
}

// GetDBSimpeg mengembalikan instance database SIMPEG
func GetDBSimpeg() *gorm.DB {
	return DBSimpeg
}

// CloseDB menutup semua koneksi database
func CloseDB() error {
	var err error
	
	// Close main DB
	if DB != nil {
		sqlDB, _ := DB.DB()
		if e := sqlDB.Close(); e != nil {
			err = e
		}
	}
	
	// Close NEOMAA
	if DBNeomaa != nil {
		sqlDB, _ := DBNeomaa.DB()
		if e := sqlDB.Close(); e != nil {
			err = e
		}
	}
	
	// Close NEOMAAREF
	if DBNeomaaRef != nil {
		sqlDB, _ := DBNeomaaRef.DB()
		if e := sqlDB.Close(); e != nil {
			err = e
		}
	}
	
	// Close SIMPEG
	if DBSimpeg != nil {
		sqlDB, _ := DBSimpeg.DB()
		if e := sqlDB.Close(); e != nil {
			err = e
		}
	}
	
	log.Println("❌ All database connections closed")
	
	return err
}