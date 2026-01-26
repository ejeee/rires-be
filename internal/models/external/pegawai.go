package external

// Pegawai represents data pegawai/dosen from NEWSIMPEG database (table: pegawai)
// Only essential fields for PKM system
type Pegawai struct {
	ID               int    `gorm:"column:id;primaryKey" json:"id"`
	NamaPegawai      string `gorm:"column:nama_pegawai" json:"nama_pegawai"`
	GelarDepan       string `gorm:"column:gelar_depan" json:"gelar_depan"`
	GelarBelakang    string `gorm:"column:gelar_belakang" json:"gelar_belakang"`
	HP               string `gorm:"column:hp" json:"hp"`
	EmailUMM         string `gorm:"column:email_umm" json:"email_umm"`
	Email            string `gorm:"column:email" json:"email"`
	IDF              int    `gorm:"column:idf" json:"idf"`                               // ID Fakultas
	HomeBaseKaryawan int    `gorm:"column:home_base_karyawan" json:"home_base_karyawan"` // ID Prodi/Unit
	IDRefAktivasi    string `gorm:"column:id_ref_aktivasi" json:"id_ref_aktivasi"`
	Hapus            int    `gorm:"column:hapus" json:"hapus"` // 0=active, 1=deleted
}

// TableName specifies the table name in SIMPEG database
func (Pegawai) TableName() string {
	return "pegawai"
}

// GetNamaLengkap returns full name with gelar
func (p *Pegawai) GetNamaLengkap() string {
	nama := ""
	if p.GelarDepan != "" {
		nama = p.GelarDepan + " "
	}
	nama += p.NamaPegawai
	if p.GelarBelakang != "" {
		nama += ", " + p.GelarBelakang
	}
	return nama
}
