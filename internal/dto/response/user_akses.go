package response

import "time"

// UserAksesResponse represents user access information
type UserAksesResponse struct {
	ID          int        `json:"id"`
	IDUserLevel int        `json:"id_user_level"`
	IDMenu      int        `json:"id_menu"`
	CanCreate   int        `json:"can_create"`
	CanUpdate   int        `json:"can_update"`
	CanDelete   int        `json:"can_delete"`
	Status      int        `json:"status"`
	TglInsert   *time.Time `json:"tgl_insert"`

	// Relations (optional)
	UserLevel *UserLevelResponse `json:"user_level,omitempty"`
	Menu      *MenuSimpleResponse      `json:"menu,omitempty"`
}

// MenuSimpleResponse represents simplified menu info
type MenuSimpleResponse struct {
	ID       int    `json:"id"`
	NamaMenu string `json:"nama_menu"`
	URLMenu  string `json:"url_menu"`
	Lucide     string `json:"lucide"`
}

// UserAksesGroupedResponse represents access grouped by user level
type UserAksesGroupedResponse struct {
	IDUserLevel int                 `json:"id_user_level"`
	NamaLevel   string              `json:"nama_level"`
	Accesses    []UserAksesResponse `json:"accesses"`
}