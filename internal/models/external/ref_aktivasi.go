package external

// RefAktivasi represents data from NEWSIMPEG table: ref_aktivasi
type RefAktivasi struct {
	ID         int    `gorm:"column:id;primaryKey" json:"id"`
	KDKODTBKOD string `gorm:"column:KDKODTBKOD" json:"kd_kod_tb_kod"`
	NMKODTBKOD string `gorm:"column:NMKODTBKOD" json:"nm_kod_tb_kod"`
	StatusPeg  int    `gorm:"column:status_peg" json:"status_peg"`
	Hapus      int    `gorm:"column:hapus" json:"hapus"`
}

// TableName specifies the table name in SIMPEG database
func (RefAktivasi) TableName() string {
	return "ref_aktivasi"
}
