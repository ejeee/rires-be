package response

import "time"

// PengajuanResponse represents detailed pengajuan data with all relations
type PengajuanResponse struct {
	ID            int                `json:"id"`
	KodePengajuan string             `json:"kode_pengajuan"`
	Judul         string             `json:"judul"`
	Tahun         int                `json:"tahun"`

	// Status
	StatusJudul    string `json:"status_judul"`
	StatusProposal string `json:"status_proposal"`
	StatusFinal    string `json:"status_final"`

	// Kategori
	Kategori *KategoriResponse `json:"kategori,omitempty"`

	// Team
	Ketua   *MahasiswaResponse  `json:"ketua,omitempty"`
	Anggota []MahasiswaResponse `json:"anggota,omitempty"`

	// Parameters (form answers)
	Parameter []ParameterResponse `json:"parameter,omitempty"`

	// File
	FileProposal    string `json:"file_proposal,omitempty"`
	FileProposalURL string `json:"file_proposal_url,omitempty"`

	// Review Judul
	ReviewerJudul      *PegawaiResponse `json:"reviewer_judul,omitempty"`
	CatatanReviewJudul string           `json:"catatan_review_judul,omitempty"`
	TglReviewJudul     *time.Time       `json:"tgl_review_judul,omitempty"`

	// Review Proposal
	ReviewerProposal      *PegawaiResponse `json:"reviewer_proposal,omitempty"`
	CatatanReviewProposal string           `json:"catatan_review_proposal,omitempty"`
	TglReviewProposal     *time.Time       `json:"tgl_review_proposal,omitempty"`

	// Review History
	ReviewJudulHistory    []ReviewResponse `json:"review_judul_history,omitempty"`
	ReviewProposalHistory []ReviewResponse `json:"review_proposal_history,omitempty"`

	// Timestamps
	TglInsert *time.Time `json:"tgl_insert"`
	TglUpdate time.Time  `json:"tgl_update"`
}

// PengajuanListResponse represents simplified pengajuan data for list view
type PengajuanListResponse struct {
	ID             int                `json:"id"`
	KodePengajuan  string             `json:"kode_pengajuan"`
	Judul          string             `json:"judul"`
	Tahun          int                `json:"tahun"`
	StatusJudul    string             `json:"status_judul"`
	StatusProposal string             `json:"status_proposal"`
	StatusFinal    string             `json:"status_final"`
	Kategori       *KategoriResponse  `json:"kategori,omitempty"`
	Ketua          *MahasiswaResponse `json:"ketua,omitempty"`
	JumlahAnggota  int                `json:"jumlah_anggota"`
	TglInsert      *time.Time         `json:"tgl_insert"`
}

// KategoriResponse represents kategori PKM data
type KategoriResponse struct {
	ID           int    `json:"id"`
	NamaKategori string `json:"nama_kategori"`
	Deskripsi    string `json:"deskripsi,omitempty"`
}

// ParameterResponse represents form parameter answer
type ParameterResponse struct {
	ID          int    `json:"id"`
	IDParameter int    `json:"id_parameter"`
	Label       string `json:"label"`
	Nilai       string `json:"nilai"`
	TipeInput   string `json:"tipe_input"`
}