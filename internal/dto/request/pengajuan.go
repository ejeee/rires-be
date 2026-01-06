package request

import "errors"

// CreatePengajuanRequest represents request body for creating PKM title submission
type CreatePengajuanRequest struct {
	IDKategori int                `json:"id_kategori" validate:"required"`
	Judul      string             `json:"judul" validate:"required,min=10,max=500"`
	Anggota    []AnggotaRequest   `json:"anggota" validate:"required,min=1,max=5,dive"`
	Parameter  []ParameterRequest `json:"parameter" validate:"dive"`
}

// UpdateJudulRequest represents request body for revising PKM title
type UpdateJudulRequest struct {
	Judul     string             `json:"judul" validate:"required,min=10,max=500"`
	Parameter []ParameterRequest `json:"parameter" validate:"dive"`
}

// Validate validates CreatePengajuanRequest
func (r *CreatePengajuanRequest) Validate() error {
	// Check if anggota has exactly 1 ketua
	ketuaCount := 0
	for _, anggota := range r.Anggota {
		if anggota.IsKetua == 1 {
			ketuaCount++
		}
	}
	
	if ketuaCount == 0 {
		return errors.New("tim harus memiliki 1 ketua")
	}
	
	if ketuaCount > 1 {
		return errors.New("tim hanya boleh memiliki 1 ketua")
	}
	
	return nil
}