package external

// Mahasiswa represents data mahasiswa from NEOMAA database
// Table structure should match the actual NEOMAA database
type Mahasiswa struct {
	NIM       string `gorm:"column:nim;primaryKey" json:"nim"`
	Nama      string `gorm:"column:nama" json:"nama"`
	Email     string `gorm:"column:email" json:"email"`
	IDProdi   int    `gorm:"column:id_prodi" json:"id_prodi"`
	Angkatan  string `gorm:"column:angkatan" json:"angkatan"`
	NoHP      string `gorm:"column:no_hp" json:"no_hp"`
	Status    int    `gorm:"column:status" json:"status"` // 1=aktif, 2=tidak aktif, dll
	
	// Relations (optional, jika perlu join)
	Prodi     *Prodi `gorm:"foreignKey:IDProdi" json:"prodi,omitempty"`
}

// TableName specifies the table name in NEOMAA database
// TODO: Sesuaikan dengan nama tabel actual di NEOMAA
func (Mahasiswa) TableName() string {
	return "master_siswa" // Change this to actual table name
}