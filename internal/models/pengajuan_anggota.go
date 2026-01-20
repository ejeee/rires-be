package models

// PengajuanAnggota represents db_pengajuan_anggota table
type PengajuanAnggota struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDPengajuan int    `gorm:"column:id_pengajuan;type:int" json:"id_pengajuan"`
	NIMAnggota  string `gorm:"column:nim_anggota;type:varchar(20)" json:"nim_anggota"` // NIM dari NEOMAA
	NamaAnggota string `gorm:"column:nama_anggota;type:varchar(100)" json:"nama_anggota"`
	IsKetua     int    `gorm:"column:is_ketua;type:int(1);default:0" json:"is_ketua"` // 1=ketua, 0=anggota
	Urutan      int    `gorm:"column:urutan;type:int" json:"urutan"`                  // 1-5
	Status      int    `gorm:"column:status;type:int(1);default:1" json:"status"`
	Hapus       int    `gorm:"column:hapus;type:int(1);default:0" json:"-"`

	// Relations
	Pengajuan *Pengajuan `gorm:"foreignKey:IDPengajuan" json:"-"`
}

// TableName specifies the table name for PengajuanAnggota model
func (PengajuanAnggota) TableName() string {
	return "db_pengajuan_anggota"
}
