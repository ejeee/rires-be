package response

import "time"

// KategoriPKMResponse untuk response kategori PKM
type KategoriPKMResponse struct {
	ID           int        `json:"id"`
	NamaKategori string     `json:"nama_kategori"`
	Status       int        `json:"status"`
	StatusText   string     `json:"status_text"` // "Aktif" atau "Tidak Aktif"
	TglInsert    *time.Time `json:"tgl_insert"`
	TglUpdate    time.Time  `json:"tgl_update"`
	UserUpdate   string     `json:"user_update"`
}

// KategoriPKMListResponse untuk response list dengan pagination
type KategoriPKMListResponse struct {
	Data       []KategoriPKMResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PerPage    int                   `json:"per_page"`
	TotalPages int                   `json:"total_pages"`
}