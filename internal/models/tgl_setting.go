package models

import "time"

// TglSetting represents db_tgl_setting table
type TglSetting struct {
	ID            int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TglDaftarAwal time.Time  `gorm:"column:tgl_daftar_awal;type:date;not null" json:"tgl_daftar_awal"`
	TglDaftarAkhir time.Time `gorm:"column:tgl_daftar_akhir;type:date;not null" json:"tgl_daftar_akhir"`
	TglReviewAwal time.Time  `gorm:"column:tgl_review_awal;type:date" json:"tgl_review_awal"`
	TglReviewAkhir time.Time `gorm:"column:tgl_review_akhir;type:date" json:"tgl_review_akhir"`
	TglPengumuman time.Time  `gorm:"column:tgl_pengumuman;type:date" json:"tgl_pengumuman"`
	Keterangan    string     `gorm:"column:keterangan;type:text" json:"keterangan"`
	IsActive      int        `gorm:"column:is_active;type:int(1);default:1" json:"is_active"` // 1=active (sedang berlaku), 0=inactive
	Status        int        `gorm:"column:status;type:int(1);default:1" json:"status"`        // 1=aktif, 2=nonaktif
	Hapus         int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`              // 0=exist, 1=deleted
	TglInsert     *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate     time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate    string     `gorm:"column:user_update;type:text" json:"user_update"`
}

// TableName specifies the table name for TglSetting model
func (TglSetting) TableName() string {
	return "db_tgl_setting"
}

// IsRegistrationOpen checks if registration is currently open
func (t *TglSetting) IsRegistrationOpen() bool {
	now := time.Now()
	return t.IsActive == 1 &&
		   t.Status == 1 &&
		   t.Hapus == 0 &&
		   now.After(t.TglDaftarAwal) &&
		   now.Before(t.TglDaftarAkhir)
}

// IsReviewPeriod checks if currently in review period
func (t *TglSetting) IsReviewPeriod() bool {
	now := time.Now()
	return t.IsActive == 1 &&
		   t.Status == 1 &&
		   t.Hapus == 0 &&
		   now.After(t.TglReviewAwal) &&
		   now.Before(t.TglReviewAkhir)
}

// IsAfterAnnouncement checks if announcement has been made
func (t *TglSetting) IsAfterAnnouncement() bool {
	now := time.Now()
	return t.IsActive == 1 &&
		   t.Status == 1 &&
		   t.Hapus == 0 &&
		   now.After(t.TglPengumuman)
}