package models

import "time"

// PlottingReviewer represents db_plotting_reviewer table
type PlottingReviewer struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDPengajuan int        `gorm:"column:id_pengajuan;type:int" json:"id_pengajuan"`
	IDPegawai   int        `gorm:"column:id_pegawai;type:int" json:"id_pegawai"`
	Tipe        string     `gorm:"column:tipe;type:varchar(20)" json:"tipe"`           // JUDUL atau PROPOSAL
	Status      string     `gorm:"column:status;type:varchar(20);default:ASSIGNED" json:"status"` // ASSIGNED, REVIEWED
	TglAssign   *time.Time `gorm:"column:tgl_assign;type:datetime" json:"tgl_assign"`
	TglInsert   *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	
	// Relations
	Pengajuan   *Pengajuan `gorm:"foreignKey:IDPengajuan" json:"-"`
}

// TableName specifies the table name for PlottingReviewer model
func (PlottingReviewer) TableName() string {
	return "db_plotting_reviewer"
}