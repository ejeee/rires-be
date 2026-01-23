package models

import "time"

// ReviewJudul represents db_review_judul table
type ReviewJudul struct {
	ID             int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDPengajuan    int        `gorm:"column:id_pengajuan;type:int" json:"id_pengajuan"`
	IDPegawai      int        `gorm:"column:id_pegawai;type:int" json:"id_pegawai"`
	IDStatusReview int        `gorm:"column:id_status_review;type:int" json:"id_status_review"` // FK ke db_status_review
	Catatan        string     `gorm:"column:catatan;type:text" json:"catatan"`
	TglReview      *time.Time `gorm:"column:tgl_review;type:datetime" json:"tgl_review"`
	Hapus          int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`
	TglInsert      *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate      time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate     string     `gorm:"column:user_update;type:text" json:"user_update"`

	// Relations
	Pengajuan    *Pengajuan    `gorm:"foreignKey:IDPengajuan" json:"-"`
	StatusReview *StatusReview `gorm:"foreignKey:IDStatusReview" json:"status_review,omitempty"`
}

// TableName specifies the table name for ReviewJudul model
func (ReviewJudul) TableName() string {
	return "db_review_judul"
}
