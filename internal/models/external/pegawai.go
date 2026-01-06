package external

// Pegawai represents data pegawai/dosen from NEWSIMPEG database
// Table structure should match the actual SIMPEG database
type Pegawai struct {
	ID              int    `gorm:"column:id;primaryKey" json:"id"`
	NIP             string `gorm:"column:nip" json:"nip"`
	Nama            string `gorm:"column:nama" json:"nama"`
	Email           string `gorm:"column:email" json:"email"`
	IDFakultas      int    `gorm:"column:id_fakultas" json:"id_fakultas"`
	BidangKeahlian  string `gorm:"column:bidang_keahlian" json:"bidang_keahlian"`
	NoHP            string `gorm:"column:no_hp" json:"no_hp"`
	Status          int    `gorm:"column:status" json:"status"`
	
	// Relations (optional)
	Fakultas        *Fakultas `gorm:"foreignKey:IDFakultas" json:"fakultas,omitempty"`
}

// TableName specifies the table name in SIMPEG database
// TODO: Sesuaikan dengan nama tabel actual di SIMPEG
func (Pegawai) TableName() string {
	return "newsimpeg" // Change this to actual table name
}