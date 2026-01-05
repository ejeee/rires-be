package models

import "time"

// User represents db_user table
type User struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NamaUser   string     `gorm:"column:nama_user;type:varchar(100)" json:"nama_user"`
	Username   string     `gorm:"column:username;type:varchar(100)" json:"username"`
	Password   string     `gorm:"column:password;type:text" json:"-"` // Hidden from JSON
	LevelUser  int        `gorm:"column:level_user;type:int" json:"level_user"`
	Status     int        `gorm:"column:status;type:int(1);default:1" json:"status"` // 1=active, 2=inactive
	Hapus      int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`       // 0=exists, 1=deleted
	TglInsert  *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate  time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate string     `gorm:"column:user_update;type:text" json:"user_update"`

	// Relations
	Level *UserLevel `gorm:"foreignKey:LevelUser" json:"level,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "db_user"
}