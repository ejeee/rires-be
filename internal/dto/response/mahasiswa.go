package response

// MahasiswaResponse represents simplified mahasiswa data for API response
type MahasiswaResponse struct {
	NIM      string         `json:"nim"`
	Nama     string         `json:"nama"`
	HP       string         `json:"no_hp"`
	Angkatan int            `json:"angkatan"`
	Prodi    *ProdiResponse `json:"program_studi,omitempty"`
	IsKetua  *int           `json:"is_ketua,omitempty"` // Only present in team context
}

// ProdiResponse represents simplified prodi data for API response
type ProdiResponse struct {
	Kode        string            `json:"kode"`
	NamaProdi   string            `json:"nama_prodi"`
	NamaSingkat string            `json:"nama_singkat"`
	Jenjang     string            `json:"jenjang"`
	Fakultas    *FakultasResponse `json:"fakultas,omitempty"`
}

// FakultasResponse represents simplified fakultas data for API response
type FakultasResponse struct {
	Kode         string `json:"kode"`
	NamaFakultas string `json:"nama_fakultas"`
	NamaSingkat  string `json:"nama_singkat"`
}