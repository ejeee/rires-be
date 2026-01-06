package external

// Prodi represents data prodi from NEOMAAREF database (table: in_programstudi)
type Prodi struct {
	Kode            string `gorm:"column:kode;primaryKey" json:"kode"`
	KodeFakultas    int    `gorm:"column:kodeFakultas" json:"kode_fakultas"`
	KodeDepart      string `gorm:"column:kode_depart" json:"kode_depart"`
	NamaDepart      string `gorm:"column:nama_depart" json:"nama_depart"`
	NamaSingkat     string `gorm:"column:nama_singkat" json:"nama_singkat"`
	Hapus           int    `gorm:"column:hapus" json:"hapus"`
	TglInsert       string `gorm:"column:tgl_insert" json:"tgl_insert"`
	TglUpdate       string `gorm:"column:tgl_update" json:"tgl_update"`
	
	// Relations
	Fakultas *Fakultas `gorm:"foreignKey:KodeFakultas;references:Kode" json:"fakultas,omitempty"`
}

// TableName specifies the table name in NEOMAAREF database
func (Prodi) TableName() string {
	return "in_programstudi"
}

// GetNamaProdi returns nama prodi (nama_depart)
func (p *Prodi) GetNamaProdi() string {
	return p.NamaDepart
}