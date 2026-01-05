package response

import "time"

// MenuResponse untuk response menu
type MenuResponse struct {
	ID         int        `json:"id"`
	ParentID   int        `json:"parent_id"`
	NamaMenu   string     `json:"nama_menu"`
	URLMenu    string     `json:"url_menu"`
	Lucide     string     `json:"lucide"`
	Urutan     int        `json:"urutan"`
	Status     int        `json:"status"`
	StatusText string     `json:"status_text"` // "Aktif" atau "Tidak Aktif"
	TglInsert  *time.Time `json:"tgl_insert"`
	TglUpdate  time.Time  `json:"tgl_update"`
	UserUpdate string     `json:"user_update"`
}

// MenuTreeResponse untuk response menu dengan children (tree structure)
type MenuTreeResponse struct {
	ID         int                `json:"id"`
	ParentID   int                `json:"parent_id"`
	NamaMenu   string             `json:"nama_menu"`
	URLMenu    string             `json:"url_menu"`
	Lucide     string             `json:"lucide"`
	Urutan     int                `json:"urutan"`
	Status     int                `json:"status"`
	StatusText string             `json:"status_text"`
	Children   []MenuTreeResponse `json:"children,omitempty"`
}

// MenuListResponse untuk response list dengan pagination
type MenuListResponse struct {
	Data       []MenuResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}