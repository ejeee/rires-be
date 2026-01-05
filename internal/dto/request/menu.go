package request

// CreateMenuRequest untuk create menu
type CreateMenuRequest struct {
	ParentID int    `json:"parent_id"` // 0 untuk main menu, >0 untuk submenu
	NamaMenu string `json:"nama_menu" validate:"required,min=3,max=50"`
	URLMenu  string `json:"url_menu" validate:"required,max=50"`
	Lucide   string `json:"lucide" validate:"max=50"` // Icon name
	Urutan   int    `json:"urutan" validate:"min=0"`
	Status   int    `json:"status" validate:"required,oneof=1 2"` // 1=active, 2=inactive
}

// UpdateMenuRequest untuk update menu
type UpdateMenuRequest struct {
	ParentID int    `json:"parent_id"`
	NamaMenu string `json:"nama_menu" validate:"required,min=3,max=50"`
	URLMenu  string `json:"url_menu" validate:"required,max=50"`
	Lucide   string `json:"lucide" validate:"max=50"`
	Urutan   int    `json:"urutan" validate:"min=0"`
	Status   int    `json:"status" validate:"required,oneof=1 2"`
}

// ReorderMenuRequest untuk ubah urutan menu
type ReorderMenuRequest struct {
	MenuID int `json:"menu_id" validate:"required"`
	Urutan int `json:"urutan" validate:"required,min=0"`
}