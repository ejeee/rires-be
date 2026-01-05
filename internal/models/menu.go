package models

import "time"

// Menu represents db_menu table
type Menu struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID   int        `gorm:"column:parent_id;type:int;default:0" json:"parent_id"`
	NamaMenu   string     `gorm:"column:nama_menu;type:varchar(50)" json:"nama_menu"`
	URLMenu    string     `gorm:"column:url_menu;type:varchar(50)" json:"url_menu"`
	Lucide     string     `gorm:"column:lucide;type:varchar(50)" json:"lucide"` // Icon name
	Urutan     int        `gorm:"column:urutan;type:int" json:"urutan"`
	Status     int        `gorm:"column:status;type:int(1);default:1" json:"status"` // 1=active, 2=inactive
	Hapus      int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`       // 0=exists, 1=deleted
	TglInsert  *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate  time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate string     `gorm:"column:user_update;type:text" json:"user_update"`
}

// TableName specifies the table name for Menu model
func (Menu) TableName() string {
	return "db_menu"
}

// MenuTree represents menu with children (for tree structure)
type MenuTree struct {
	Menu
	Children []MenuTree `json:"children,omitempty"`
}