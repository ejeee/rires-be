package request

// CreateStatusReviewRequest untuk create status review
type CreateStatusReviewRequest struct {
	NamaStatus string `json:"nama_status" validate:"required,min=3,max=50"`
	KodeStatus string `json:"kode_status" validate:"required,min=2,max=20,uppercase"`
	Warna      string `json:"warna" validate:"required,max=20"` // gray, blue, green, yellow, red, etc
	Urutan     int    `json:"urutan" validate:"min=0"`
	Status     int    `json:"status" validate:"required,oneof=1 2"` // 1=active, 2=inactive
}

// UpdateStatusReviewRequest untuk update status review
type UpdateStatusReviewRequest struct {
	NamaStatus string `json:"nama_status" validate:"required,min=3,max=50"`
	KodeStatus string `json:"kode_status" validate:"required,min=2,max=20,uppercase"`
	Warna      string `json:"warna" validate:"required,max=20"`
	Urutan     int    `json:"urutan" validate:"min=0"`
	Status     int    `json:"status" validate:"required,oneof=1 2"`
}