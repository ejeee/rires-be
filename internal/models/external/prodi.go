package external

// Prodi represents data prodi from NEOMAAREF database
type Prodi struct {
	ID          int    `gorm:"column:id;primaryKey" json:"id"`
	KodeProdi   string `gorm:"column:kode_prodi" json:"kode_prodi"`
	NamaProdi   string `gorm:"column:nama_prodi" json:"nama_prodi"`
	IDFakultas  int    `gorm:"column:id_fakultas" json:"id_fakultas"`
	Jenjang     string `gorm:"column:jenjang" json:"jenjang"` // S1, S2, S3, D3, D4
	Status      int    `gorm:"column:status" json:"status"`
	
	// Relations
	Fakultas    *Fakultas `gorm:"foreignKey:IDFakultas" json:"fakultas,omitempty"`
}

// TableName specifies the table name in NEOMAAREF database
// TODO: Sesuaikan dengan nama tabel actual di NEOMAAREF
func (Prodi) TableName() string {
	return "in_programstudi" // Change this to actual table name
}