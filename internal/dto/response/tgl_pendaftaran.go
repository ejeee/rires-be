package response

import "time"

// TanggalPendaftaranResponse untuk response tanggal pendaftaran
type TanggalPendaftaranResponse struct {
	ID              int       `json:"id"`
	TanggalMulai    time.Time `json:"tanggal_mulai"`
	TanggalSelesai  time.Time `json:"tanggal_selesai"`
	Keterangan      string    `json:"keterangan"`
	IsActive        int       `json:"is_active"`
	IsActiveText    string    `json:"is_active_text"` // "Aktif" atau "Tidak Aktif"
	IsOpen          bool      `json:"is_open"`        // Currently open for registration?
	Status          int       `json:"status"`
	StatusText      string    `json:"status_text"`     // "Aktif" atau "Tidak Aktif"
	DaysRemaining   int       `json:"days_remaining"`  // Hari tersisa
	TglInsert       *time.Time `json:"tgl_insert"`
	TglUpdate       time.Time `json:"tgl_update"`
	UserUpdate      string    `json:"user_update"`
}

// TanggalPendaftaranListResponse untuk response list dengan pagination
type TanggalPendaftaranListResponse struct {
	Data       []TanggalPendaftaranResponse `json:"data"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	PerPage    int                          `json:"per_page"`
	TotalPages int                          `json:"total_pages"`
}

// RegistrationStatusResponse untuk cek status pendaftaran (public)
type RegistrationStatusResponse struct {
	IsOpen         bool      `json:"is_open"`
	Message        string    `json:"message"`
	TanggalMulai   time.Time `json:"tanggal_mulai,omitempty"`
	TanggalSelesai time.Time `json:"tanggal_selesai,omitempty"`
	DaysRemaining  int       `json:"days_remaining,omitempty"`
}