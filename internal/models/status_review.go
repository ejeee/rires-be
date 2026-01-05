package models

import "time"

// StatusReview represents db_status_review table
type StatusReview struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NamaStatus string     `gorm:"column:nama_status;type:varchar(50)" json:"nama_status"`
	KodeStatus string     `gorm:"column:kode_status;type:varchar(20);uniqueIndex" json:"kode_status"`
	Warna      string     `gorm:"column:warna;type:varchar(20)" json:"warna"` // For UI badge color
	Urutan     int        `gorm:"column:urutan;type:int;default:0" json:"urutan"`
	Status     int        `gorm:"column:status;type:int(1);default:1" json:"status"` // 1=active, 2=inactive
	Hapus      int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`       // 0=exists, 1=deleted
	TglInsert  *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate  time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate string     `gorm:"column:user_update;type:text" json:"user_update"`
}

// TableName specifies the table name for StatusReview model
func (StatusReview) TableName() string {
	return "db_status_review"
}