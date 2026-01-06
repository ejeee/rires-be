package request

// CreateTanggalPendaftaranRequest untuk create tanggal pendaftaran
type CreateTanggalPendaftaranRequest struct {
	TanggalMulai   string `json:"tanggal_mulai" validate:"required"`   // Format: "2026-01-01 00:00:00"
	TanggalSelesai string `json:"tanggal_selesai" validate:"required"` // Format: "2026-01-31 23:59:59"
	Keterangan     string `json:"keterangan"`
	Status         int    `json:"status" validate:"required,oneof=1 2"` // 1=enabled, 2=disabled
}

// UpdateTanggalPendaftaranRequest untuk update tanggal pendaftaran
type UpdateTanggalPendaftaranRequest struct {
	TanggalMulai   string `json:"tanggal_mulai" validate:"required"`
	TanggalSelesai string `json:"tanggal_selesai" validate:"required"`
	Keterangan     string `json:"keterangan"`
	Status         int    `json:"status" validate:"required,oneof=1 2"`
}