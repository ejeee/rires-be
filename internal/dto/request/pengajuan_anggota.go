package request

// AnggotaRequest represents team member data in pengajuan request
type AnggotaRequest struct {
	NIM     string `json:"nim" validate:"required"`
	IsKetua int    `json:"is_ketua" validate:"oneof=0 1"`           // 1=ketua, 0=anggota
	Urutan  int    `json:"urutan" validate:"omitempty,min=1,max=5"` // Optional, will be auto-assigned
}
