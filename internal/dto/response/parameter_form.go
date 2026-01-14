package response

import "time"

// ParameterFormResponse untuk response parameter form
type ParameterFormResponse struct {
	ID            int        `json:"id"`
	IDKategori    int        `json:"id_kategori"`
	NamaKategori  string     `json:"nama_kategori,omitempty"` // Dari join
	NamaParameter string     `json:"nama_parameter"`
	Label         string     `json:"label"`
	TipeInput     string     `json:"tipe_input"`
	Validasi      string     `json:"validasi"` // JSON string
	Placeholder   string     `json:"placeholder"`
	HelpText      string     `json:"help_text"`
	Opsi          string     `json:"opsi"` // JSON string for radio/select options
	Urutan        int        `json:"urutan"`
	Status        int        `json:"status"`
	StatusText    string     `json:"status_text"` // "Aktif" atau "Tidak Aktif"
	TglInsert     *time.Time `json:"tgl_insert"`
	TglUpdate     time.Time  `json:"tgl_update"`
	UserUpdate    string     `json:"user_update"`
}

// ParameterFormListResponse untuk response list dengan pagination
type ParameterFormListResponse struct {
	Data       []ParameterFormResponse `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
	TotalPages int                     `json:"total_pages"`
}

// ParameterFormByKategoriResponse untuk response grouped by kategori
type ParameterFormByKategoriResponse struct {
	IDKategori   int                     `json:"id_kategori"`
	NamaKategori string                  `json:"nama_kategori"`
	Parameters   []ParameterFormResponse `json:"parameters"`
}
