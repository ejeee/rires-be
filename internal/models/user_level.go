package models

import "time"

// UserLevel represents db_user_level table
type UserLevel struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NamaLevel    string     `gorm:"column:nama_level;type:varchar(100)" json:"nama_level"`
	IDConference int        `gorm:"column:id_conference;type:int;default:0" json:"id_conference"` // Add this field
	Status       int        `gorm:"column:status;type:int(1);default:1" json:"status"`             // 1=active, 2=inactive
	Hapus        int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`                   // 0=exists, 1=deleted (hidden from JSON)
	TglInsert    *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate    time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate   string        `gorm:"column:user_update;type:int" json:"user_update"`
}

// TableName specifies the table name for UserLevel model
func (UserLevel) TableName() string {
	return "db_user_level"
}