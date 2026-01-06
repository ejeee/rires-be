package models

import "time"

// TanggalPendaftaran represents db_tanggal_pendaftaran table
type TanggalPendaftaran struct {
	ID              int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TanggalMulai    time.Time  `gorm:"column:tanggal_mulai;type:datetime" json:"tanggal_mulai"`
	TanggalSelesai  time.Time  `gorm:"column:tanggal_selesai;type:datetime" json:"tanggal_selesai"`
	Keterangan      string     `gorm:"column:keterangan;type:text" json:"keterangan"`
	IsActive        int        `gorm:"column:is_active;type:int(1);default:1" json:"is_active"` // 1=currently active period, 0=not active
	Status          int        `gorm:"column:status;type:int(1);default:1" json:"status"`        // 1=enabled, 2=disabled
	Hapus           int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`              // 0=exists, 1=deleted
	TglInsert       *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate       time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate      string     `gorm:"column:user_update;type:text" json:"user_update"`
}

// TableName specifies the table name for TanggalPendaftaran model
func (TanggalPendaftaran) TableName() string {
	return "db_tanggal_pendaftaran"
}

// IsOpen checks if registration is currently open
func (t *TanggalPendaftaran) IsOpen() bool {
	now := time.Now()
	return t.IsActive == 1 && 
		   t.Status == 1 && 
		   t.Hapus == 0 &&
		   now.After(t.TanggalMulai) && 
		   now.Before(t.TanggalSelesai)
}