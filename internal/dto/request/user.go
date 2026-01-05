package request

// LoginRequest adalah struktur untuk request login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// CreateUserRequest untuk create user
type CreateUserRequest struct {
	NamaUser  string `json:"nama_user" validate:"required,min=3,max:100"`
	Username  string `json:"username" validate:"required,min=3,max:100"`
	Password  string `json:"password" validate:"required,min=6"` // Plain password, akan di-hash
	LevelUser int    `json:"level_user" validate:"required"`
	Status    int    `json:"status" validate:"required,oneof=1 2"` // 1=active, 2=inactive
}

// UpdateUserRequest untuk update user
type UpdateUserRequest struct {
	NamaUser  string `json:"nama_user" validate:"required,min=3,max:100"`
	Username  string `json:"username" validate:"required,min=3,max:100"`
	LevelUser int    `json:"level_user" validate:"required"`
	Status    int    `json:"status" validate:"required,oneof=1 2"`
}

// ChangePasswordRequest untuk ubah password user
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ResetPasswordRequest untuk reset password (admin only)
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=6"`
}