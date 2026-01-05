package response

import "time"

// UserResponse untuk response user (without password)
type UserResponse struct {
	ID            int        `json:"id"`
	NamaUser      string     `json:"nama_user"`
	Username      string     `json:"username"`
	LevelUser     int        `json:"level_user"`
	NamaLevel     string     `json:"nama_level,omitempty"` // Dari join
	Status        int        `json:"status"`
	StatusText    string     `json:"status_text"` // "Aktif" atau "Tidak Aktif"
	TglInsert     *time.Time `json:"tgl_insert"`
	TglUpdate     time.Time  `json:"tgl_update"`
	UserUpdate    string     `json:"user_update"`
	LastLoginText string     `json:"last_login_text,omitempty"` // Optional
}

// UserListResponse untuk response list dengan pagination
type UserListResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}