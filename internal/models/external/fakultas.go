package external

// Fakultas represents data fakultas from NEOMAAREF database (table: in_fakultas)
type Fakultas struct {
	Kode          string `gorm:"column:kode;primaryKey" json:"kode"`
	NamaFakultas  string `gorm:"column:namaFakultas" json:"nama_fakultas"`
	NamaFakPendek string `gorm:"column:namaFakPendek" json:"nama_fak_pendek"`
	STAktif       int    `gorm:"column:st_aktif" json:"st_aktif"`
	Hapus         int    `gorm:"column:hapus" json:"hapus"`
	TglInsert     string `gorm:"column:tgl_insert" json:"tgl_insert"`
	TglUpdate     string `gorm:"column:tgl_update" json:"tgl_update"`
	UserUpdate    string `gorm:"column:user_update" json:"user_update"`
}

// TableName specifies the table name in NEOMAAREF database
func (Fakultas) TableName() string {
	return "in_fakultas"
}