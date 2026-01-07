package models

import "time"

// UserAkses represents db_user_akses table
type UserAkses struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDUserLevel int        `gorm:"column:id_user_level;type:int(11)" json:"id_user_level"`
	IDMenu      int        `gorm:"column:id_menu;type:int(11)" json:"id_menu"`
	CanCreate   int        `gorm:"column:can_create;type:int(1);default:0" json:"can_create"` // 1=can create, 0=cannot
	CanUpdate   int        `gorm:"column:can_update;type:int(1);default:0" json:"can_update"` // 1=can update, 0=cannot
	CanDelete   int        `gorm:"column:can_delete;type:int(1);default:0" json:"can_delete"` // 1=can delete, 0=cannot
	Status      int        `gorm:"column:status;type:int(1);default:1" json:"status"`
	Hapus       int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`
	TglInsert   *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate   time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate  string     `gorm:"column:user_update;type:text" json:"user_update"`

	// Relations (optional - for preload)
	UserLevel *UserLevel `gorm:"foreignKey:IDUserLevel" json:"user_level,omitempty"`
	Menu      *Menu      `gorm:"foreignKey:IDMenu" json:"menu,omitempty"`
}

// TableName specifies the table name for UserAkses model
func (UserAkses) TableName() string {
	return "db_user_akses"
}

// HasCreatePermission checks if user has create permission
func (u *UserAkses) HasCreatePermission() bool {
	return u.CanCreate == 1 && u.Status == 1 && u.Hapus == 0
}

// HasUpdatePermission checks if user has update permission
func (u *UserAkses) HasUpdatePermission() bool {
	return u.CanUpdate == 1 && u.Status == 1 && u.Hapus == 0
}

// HasDeletePermission checks if user has delete permission
func (u *UserAkses) HasDeletePermission() bool {
	return u.CanDelete == 1 && u.Status == 1 && u.Hapus == 0
}