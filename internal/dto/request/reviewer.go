package request

// ActivateReviewerRequest represents request to activate pegawai as reviewer
type ActivateReviewerRequest struct {
	IDPegawai int `json:"id_pegawai" validate:"required"`
}

// UpdateReviewerRequest represents request to update reviewer data
type UpdateReviewerRequest struct {
	IsActive int `json:"is_active" validate:"oneof=0 1"`
}