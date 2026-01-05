package models

import (
	"time"
)

// User adalah model untuk tabel db_user (existing database)
type User struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NamaUser   string    `gorm:"column:nama_user;type:varchar(100);not null" json:"nama_user"`
	Username   string    `gorm:"column:username;type:varchar(100);not null" json:"username"`
	Password   string    `gorm:"column:password;type:text;not null" json:"-"` // json:"-" agar tidak muncul di response
	LevelUser  int       `gorm:"column:level_user;type:int(11);not null" json:"level_user"`
	Status     int       `gorm:"column:status;type:int(1);not null;default:1" json:"status"` // 1: aktif, 2: tidak aktif
	Hapus      int       `gorm:"column:hapus;type:int(1);not null;default:0" json:"-"`       // 0: ada, 1: hapus (soft delete)
	TglInsert  time.Time `gorm:"column:tgl_insert;type:datetime;not null" json:"tgl_insert"`
	TglUpdate  time.Time `gorm:"column:tgl_update;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"tgl_update"`
	UserUpdate string    `gorm:"column:user_update;type:text;not null" json:"user_update"`
}

// TableName menentukan nama tabel di database
func (User) TableName() string {
	return "db_user"
}