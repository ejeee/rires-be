package models

import "time"

// KategoriPKM represents db_kategori_pkm table
type KategoriPKM struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NamaKategori string     `gorm:"column:nama_kategori;type:varchar(100)" json:"nama_kategori"`
	Status       int        `gorm:"column:status;type:int(1);default:1" json:"status"` // 1=active, 2=inactive
	Hapus        int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`       // 0=exists, 1=deleted
	TglInsert    *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate    time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate   string     `gorm:"column:user_update;type:text" json:"user_update"`
}

// TableName specifies the table name for KategoriPKM model
func (KategoriPKM) TableName() string {
	return "db_kategori_pkm"
}