package response

import "time"

// StatusReviewResponse untuk response status review
type StatusReviewResponse struct {
	ID         int        `json:"id"`
	NamaStatus string     `json:"nama_status"`
	KodeStatus string     `json:"kode_status"`
	Warna      string     `json:"warna"`
	Urutan     int        `json:"urutan"`
	Status     int        `json:"status"`
	StatusText string     `json:"status_text"` // "Aktif" atau "Tidak Aktif"
	TglInsert  *time.Time `json:"tgl_insert"`
	TglUpdate  time.Time  `json:"tgl_update"`
	UserUpdate string     `json:"user_update"`
}

// StatusReviewListResponse untuk response list dengan pagination
type StatusReviewListResponse struct {
	Data       []StatusReviewResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PerPage    int                    `json:"per_page"`
	TotalPages int                    `json:"total_pages"`
}