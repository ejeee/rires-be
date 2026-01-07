package request

// CreateUserAksesRequest represents request to create user access
type CreateUserAksesRequest struct {
	IDUserLevel int `json:"id_user_level" validate:"required"`
	IDMenu      int `json:"id_menu" validate:"required"`
	CanCreate   int `json:"can_create" validate:"oneof=0 1"`
	CanUpdate   int `json:"can_update" validate:"oneof=0 1"`
	CanDelete   int `json:"can_delete" validate:"oneof=0 1"`
}

// UpdateUserAksesRequest represents request to update user access
type UpdateUserAksesRequest struct {
	CanCreate int `json:"can_create" validate:"oneof=0 1"`
	CanUpdate int `json:"can_update" validate:"oneof=0 1"`
	CanDelete int `json:"can_delete" validate:"oneof=0 1"`
}

// BulkCreateUserAksesRequest represents request to create multiple accesses at once
type BulkCreateUserAksesRequest struct {
	IDUserLevel int                     `json:"id_user_level" validate:"required"`
	Menus       []MenuPermissionRequest `json:"menus" validate:"required,dive"`
}

// MenuPermissionRequest represents permission for one menu
type MenuPermissionRequest struct {
	IDMenu    int `json:"id_menu" validate:"required"`
	CanCreate int `json:"can_create" validate:"oneof=0 1"`
	CanUpdate int `json:"can_update" validate:"oneof=0 1"`
	CanDelete int `json:"can_delete" validate:"oneof=0 1"`
}