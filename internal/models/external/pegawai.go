package external

// Pegawai represents data pegawai/dosen from NEWSIMPEG database (table: pegawai)
type Pegawai struct {
	ID                    int    `gorm:"column:id;primaryKey" json:"id"`
	NIP                   string `gorm:"column:nip" json:"nip"`
	NamaPegawai           string `gorm:"column:nama_pegawai" json:"nama_pegawai"`
	GelarDepan            string `gorm:"column:gelar_depan" json:"gelar_depan"`
	GelarBelakang         string `gorm:"column:gelar_belakang" json:"gelar_belakang"`
	HP                    string `gorm:"column:hp" json:"hp"`
	Foto                  string `gorm:"column:foto" json:"foto"`
	EmailUMM              string `gorm:"column:email_umm" json:"email_umm"`
	Email                 string `gorm:"column:email" json:"email"`
	Hapus                 int    `gorm:"column:hapus" json:"hapus"` // 1=ada, 0=hapus
	TglInsert             string `gorm:"column:tgl_insert" json:"tgl_insert"`
	TglUpdate             string `gorm:"column:tgl_update" json:"tgl_update"`
	UserUpdate            string `gorm:"column:user_update" json:"user_update"`
	
	// Relations (optional)
	Fakultas *Fakultas `gorm:"foreignKey:IDF;references:Kode" json:"fakultas,omitempty"`
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