package external

// Mahasiswa represents data mahasiswa from NEOMAA database (table: master_siswa)
// Only essential fields for PKM system
type Mahasiswa struct {
	KodeSiswa       string `gorm:"column:kode_siswa;primaryKey" json:"kode_siswa"` // NIM
	NamaSiswa       string `gorm:"column:nama_siswa" json:"nama_siswa"`             // Nama lengkap
	HPSiswa         string `gorm:"column:hp_siswa" json:"hp_siswa"`                 // No HP
	RefProgramStudi int    `gorm:"column:ref_program_studi" json:"ref_program_studi"` // ID Prodi
	TahunMasuk      int    `gorm:"column:tahun_masuk" json:"tahun_masuk"`           // Angkatan
	
	// Relations (optional, jika perlu join)
	Prodi *Prodi `gorm:"foreignKey:RefProgramStudi;references:Kode" json:"prodi,omitempty"`
}

// TableName specifies the table name in NEOMAA database
func (Mahasiswa) TableName() string {
	return "master_siswa"
}

// GetNIM returns the NIM (kode_siswa)
func (m *Mahasiswa) GetNIM() string {
	return m.KodeSiswa
}

// GetNama returns the nama (nama_siswa)
func (m *Mahasiswa) GetNama() string {
	return m.NamaSiswa
}

// GetAngkatan returns the angkatan (tahun_masuk)
func (m *Mahasiswa) GetAngkatan() string {
	if m.TahunMasuk > 0 {
		return string(rune(m.TahunMasuk))
	}
	return ""
}