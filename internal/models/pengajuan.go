package models

import "time"

// Pengajuan represents db_pengajuan_pkm table
type Pengajuan struct {
	ID              int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	KodePengajuan   string     `gorm:"column:kode_pengajuan;type:varchar(50);uniqueIndex" json:"kode_pengajuan"`
	NamaKetua       string     `gorm:"column:nama_ketua;type:varchar(100)" json:"nama_ketua"`
	NIMKetua        string     `gorm:"column:nim_ketua;type:varchar(20)" json:"nim_ketua"`
	IDKategori      int        `gorm:"column:id_kategori;type:int" json:"id_kategori"`
	Judul           string     `gorm:"column:judul;type:text" json:"judul"`
	EmailKetua      string     `gorm:"column:email_ketua;type:varchar(100)" json:"email_ketua"`
	NoHPKetua       string     `gorm:"column:no_hp_ketua;type:varchar(50)" json:"no_hp_ketua"`
	ProgramStudi    string     `gorm:"column:program_studi;type:varchar(100)" json:"program_studi"`
	Fakultas        string     `gorm:"column:fakultas;type:varchar(100)" json:"fakultas"`
	DosenPembimbing string     `gorm:"column:dosen_pembimbing;type:varchar(100)" json:"dosen_pembimbing"`
	ParameterData   string     `gorm:"column:parameter_data;type:text" json:"parameter_data"` // JSON string for form parameters
	TglPengajuan    *time.Time `gorm:"column:tgl_pengajuan;type:datetime" json:"tgl_pengajuan"`
	Tahun           int        `gorm:"column:tahun;type:int" json:"tahun"`

	// Status Judul
	StatusJudul            string     `gorm:"column:status_judul;type:varchar(20);default:PENDING" json:"status_judul"` // PENDING, ON_REVIEW, ACC, REVISI, TOLAK
	IDPegawaiReviewerJudul *int       `gorm:"column:id_pegawai_reviewer_judul;type:int" json:"id_pegawai_reviewer_judul"`
	CatatanReviewJudul     string     `gorm:"column:catatan_review_judul;type:text" json:"catatan_review_judul"`
	TglReviewJudul         *time.Time `gorm:"column:tgl_review_judul;type:datetime" json:"tgl_review_judul"`

	// Status Proposal
	FileProposal              string     `gorm:"column:file_proposal;type:text" json:"file_proposal"`
	StatusProposal            string     `gorm:"column:status_proposal;type:varchar(20)" json:"status_proposal"` // PENDING, ON_REVIEW, ACC, REVISI, TOLAK
	IDPegawaiReviewerProposal *int       `gorm:"column:id_pegawai_reviewer_proposal;type:int" json:"id_pegawai_reviewer_proposal"`
	CatatanReviewProposal     string     `gorm:"column:catatan_review_proposal;type:text" json:"catatan_review_proposal"`
	TglReviewProposal         *time.Time `gorm:"column:tgl_review_proposal;type:datetime" json:"tgl_review_proposal"`

	// Final Status
	StatusFinal string `gorm:"column:status_final;type:varchar(20);default:DRAFT" json:"status_final"` // DRAFT, SUBMITTED, LOLOS, TIDAK_LOLOS

	Status     int        `gorm:"column:status;type:int(1);default:1" json:"status"`
	Hapus      int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`
	TglInsert  *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate  time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate string     `gorm:"column:user_update;type:text" json:"user_update"`

	// Relations (will be loaded via Preload)
	Kategori       *KategoriPKM       `gorm:"foreignKey:IDKategori" json:"kategori,omitempty"`
	Anggota        []PengajuanAnggota `gorm:"foreignKey:IDPengajuan" json:"anggota,omitempty"`
	Parameter      []ParameterPKM     `gorm:"foreignKey:IDPengajuan" json:"parameter,omitempty"`
	ReviewJudul    []ReviewJudul      `gorm:"foreignKey:IDPengajuan" json:"review_judul,omitempty"`
	ReviewProposal []ReviewProposal   `gorm:"foreignKey:IDPengajuan" json:"review_proposal,omitempty"`
}

// TableName specifies the table name for Pengajuan model
func (Pengajuan) TableName() string {
	return "db_pengajuan_pkm"
}

// CanUploadProposal checks if mahasiswa can upload proposal (status_judul must be ACC)
func (p *Pengajuan) CanUploadProposal() bool {
	return p.StatusJudul == "ACC"
}

// CanReviseJudul checks if mahasiswa can revise judul (status_judul must be REVISI)
func (p *Pengajuan) CanReviseJudul() bool {
	return p.StatusJudul == "REVISI"
}

// CanReviseProposal checks if mahasiswa can revise proposal (status_proposal must be REVISI)
func (p *Pengajuan) CanReviseProposal() bool {
	return p.StatusProposal == "REVISI"
}

// IsOwner checks if given NIM is the owner (ketua) of this pengajuan
func (p *Pengajuan) IsOwner(nim string) bool {
	return p.NIMKetua == nim
}
