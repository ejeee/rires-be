package request

// CreateKategoriPKMRequest untuk create kategori PKM
type CreateKategoriPKMRequest struct {
	NamaKategori string `json:"nama_kategori" validate:"required,min=3,max=100"`
	Status       int    `json:"status" validate:"required,oneof=1 2"` // 1=active, 2=inactive
}

// UpdateKategoriPKMRequest untuk update kategori PKM
type UpdateKategoriPKMRequest struct {
	NamaKategori string `json:"nama_kategori" validate:"required,min=3,max=100"`
	Status       int    `json:"status" validate:"required,oneof=1 2"`
}