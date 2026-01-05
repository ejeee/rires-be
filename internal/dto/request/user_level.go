package request

// CreateUserLevelRequest untuk create user level
type CreateUserLevelRequest struct {
	NamaLevel string `json:"nama_level" validate:"required,min=3,max=100"`
	Status    int    `json:"status" validate:"required,oneof=1 2"` // 1=active, 2=inactive
}

// UpdateUserLevelRequest untuk update user level
type UpdateUserLevelRequest struct {
	NamaLevel string `json:"nama_level" validate:"required,min=3,max=100"`
	Status    int    `json:"status" validate:"required,oneof=1 2"`
}