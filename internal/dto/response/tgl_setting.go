package response

import "time"

// TglSettingResponse untuk response tanggal setting
type TglSettingResponse struct {
	ID             int        `json:"id"`
	TglDaftarAwal  time.Time  `json:"tgl_daftar_awal"`
	TglDaftarAkhir time.Time  `json:"tgl_daftar_akhir"`
	TglReviewAwal  time.Time  `json:"tgl_review_awal"`
	TglReviewAkhir time.Time  `json:"tgl_review_akhir"`
	TglPengumuman  time.Time  `json:"tgl_pengumuman"`
	Keterangan     string     `json:"keterangan"`
	IsActive       int        `json:"is_active"`
	IsActiveText   string     `json:"is_active_text"`        // "Aktif" atau "Tidak Aktif"
	IsRegOpen      bool       `json:"is_reg_open"`           // Registration period open?
	IsReviewPeriod bool       `json:"is_review_period"`      // In review period?
	IsAnnounced    bool       `json:"is_announced"`          // After announcement?
	Status         int        `json:"status"`
	StatusText     string     `json:"status_text"`           // "Aktif" atau "Tidak Aktif"
	DaysRemaining  int        `json:"days_remaining"`        // Hari tersisa pendaftaran
	TglInsert      *time.Time `json:"tgl_insert"`
	TglUpdate      time.Time  `json:"tgl_update"`
	UserUpdate     string     `json:"user_update"`
}

// TglSettingListResponse untuk response list dengan pagination
type TglSettingListResponse struct {
	Data       []TglSettingResponse `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"per_page"`
	TotalPages int                  `json:"total_pages"`
}

// RegistrationStatusResponse untuk cek status pendaftaran (public)
type RegistrationStatusResponse struct {
	IsOpen         bool      `json:"is_open"`
	Message        string    `json:"message"`
	TglDaftarAwal  time.Time `json:"tgl_daftar_awal,omitempty"`
	TglDaftarAkhir time.Time `json:"tgl_daftar_akhir,omitempty"`
	TglReviewAwal  time.Time `json:"tgl_review_awal,omitempty"`
	TglReviewAkhir time.Time `json:"tgl_review_akhir,omitempty"`
	TglPengumuman  time.Time `json:"tgl_pengumuman,omitempty"`
	DaysRemaining  int       `json:"days_remaining,omitempty"`
}