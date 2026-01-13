package request

// CreateTglSettingRequest untuk create tanggal setting
type CreateTglSettingRequest struct {
	TglDaftarAwal  string `json:"tgl_daftar_awal" validate:"required"`   // Format: "2026-01-01"
	TglDaftarAkhir string `json:"tgl_daftar_akhir" validate:"required"`  // Format: "2026-01-31"
	TglReviewAwal  string `json:"tgl_review_awal"`                       // Format: "2026-02-01"
	TglReviewAkhir string `json:"tgl_review_akhir"`                      // Format: "2026-02-28"
	TglPengumuman  string `json:"tgl_pengumuman"`                        // Format: "2026-03-01"
	Keterangan     string `json:"keterangan"`
	Status         int    `json:"status" validate:"required,oneof=1 2"`  // 1=aktif, 2=nonaktif
}

// UpdateTglSettingRequest untuk update tanggal setting
type UpdateTglSettingRequest struct {
	TglDaftarAwal  string `json:"tgl_daftar_awal" validate:"required"`
	TglDaftarAkhir string `json:"tgl_daftar_akhir" validate:"required"`
	TglReviewAwal  string `json:"tgl_review_awal"`
	TglReviewAkhir string `json:"tgl_review_akhir"`
	TglPengumuman  string `json:"tgl_pengumuman"`
	Keterangan     string `json:"keterangan"`
	Status         int    `json:"status" validate:"required,oneof=1 2"`
}