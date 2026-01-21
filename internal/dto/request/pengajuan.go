package request

import "errors"

// CreatePengajuanRequest represents request body for creating PKM title submission
type CreatePengajuanRequest struct {
	IDKategori      int                    `json:"id_kategori" validate:"required"`
	Judul           string                 `json:"judul" validate:"required,min=10,max=500"`
	NIMKetua        string                 `json:"nim_ketua" validate:"omitempty"` // Optional: for admin to specify ketua NIM
	NamaKetua       string                 `json:"nama_ketua" validate:"omitempty"`
	EmailKetua      string                 `json:"email_ketua" validate:"omitempty,email"`
	NoHPKetua       string                 `json:"no_hp_ketua" validate:"omitempty"`
	ProgramStudi    string                 `json:"program_studi" validate:"omitempty"`
	Fakultas        string                 `json:"fakultas" validate:"omitempty"`
	DosenPembimbing string                 `json:"dosen_pembimbing" validate:"omitempty"`
	Anggota         []AnggotaRequest       `json:"anggota" validate:"omitempty,max=5,dive"` // Optional, ketua auto-added
	ParameterData   map[string]interface{} `json:"parameter_data" validate:"omitempty"`     // JSON object for form parameters
}

// UpdateJudulRequest represents request body for revising PKM title
type UpdateJudulRequest struct {
	IDKategori      int                    `json:"id_kategori" validate:"omitempty"`
	Judul           string                 `json:"judul" validate:"required,min=10,max=500"`
	NamaKetua       string                 `json:"nama_ketua" validate:"omitempty"`
	EmailKetua      string                 `json:"email_ketua" validate:"omitempty,email"`
	NoHPKetua       string                 `json:"no_hp_ketua" validate:"omitempty"`
	ProgramStudi    string                 `json:"program_studi" validate:"omitempty"`
	Fakultas        string                 `json:"fakultas" validate:"omitempty"`
	DosenPembimbing string                 `json:"dosen_pembimbing" validate:"omitempty"`
	Anggota         []AnggotaRequest       `json:"anggota" validate:"omitempty,max=5,dive"`
	ParameterData   map[string]interface{} `json:"parameter_data" validate:"omitempty"`
}

// Validate validates CreatePengajuanRequest
// Note: Ketua validation is now handled by service (auto-added based on authenticated user)
func (r *CreatePengajuanRequest) Validate() error {
	// Check if anggota has exactly 1 ketua
	ketuaCount := 0
	for _, anggota := range r.Anggota {
		if anggota.IsKetua == 1 {
			ketuaCount++
		}
	}

	// Allow 0 ketua (will be auto-added by service)
	// Only check if there are more than 1 ketua
	if ketuaCount > 1 {
		return errors.New("tim hanya boleh memiliki 1 ketua")
	}

	return nil
}
