package models

import "time"

// Reviewer represents db_reviewer table
type Reviewer struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDPegawai    int        `gorm:"column:id_pegawai;type:int(11);uniqueIndex" json:"id_pegawai"`
	NamaPegawai  string     `gorm:"column:nama_pegawai;type:varchar(255)" json:"nama_pegawai"`
	EmailUmm     string     `gorm:"column:email_umm;type:varchar(255)" json:"email_umm"`
	IsActive     int        `gorm:"column:is_active;type:int(1);default:1" json:"is_active"` // 1=active, 0=inactive
	Status       int        `gorm:"column:status;type:int(1);default:1" json:"status"`
	Hapus        int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`
	TglInsert    *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate    time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate   string     `gorm:"column:user_update;type:text" json:"user_update"`
}

// TableName specifies the table name for Reviewer model
func (Reviewer) TableName() string {
	return "db_reviewer"
}

// IsActiveReviewer checks if reviewer is active
func (r *Reviewer) IsActiveReviewer() bool {
	return r.IsActive == 1 && r.Status == 1 && r.Hapus == 0
}