package external

// Fakultas represents data fakultas from NEOMAAREF database
type Fakultas struct {
	ID           int    `gorm:"column:id;primaryKey" json:"id"`
	KodeFakultas string `gorm:"column:kode_fakultas" json:"kode_fakultas"`
	NamaFakultas string `gorm:"column:nama_fakultas" json:"nama_fakultas"`
	NamaSingkat  string `gorm:"column:nama_singkat" json:"nama_singkat"`
	Status       int    `gorm:"column:status" json:"status"`
}

// TableName specifies the table name in NEOMAAREF database
// TODO: Sesuaikan dengan nama tabel actual di NEOMAAREF
func (Fakultas) TableName() string {
	return "in_fakultas" // Change this to actual table name
}