package response

// FakultasResponse represents fakultas data for API response
type FakultasRefResponse struct {
	Kode          string `json:"kode"`
	NamaFakultas  string `json:"nama_fakultas"`
	NamaFakPendek string `json:"nama_fak_pendek"`
}

// ProdiResponse represents prodi data for API response
type ProdiRefResponse struct {
	Kode         string               `json:"kode"`
	KodeFakultas int                  `json:"kode_fakultas"`
	KodeDepart   string               `json:"kode_depart"`
	NamaDepart   string               `json:"nama_depart"`
	NamaSingkat  string               `json:"nama_singkat"`
	Fakultas     *FakultasRefResponse `json:"fakultas,omitempty"`
}

// FakultasListResponse untuk response list fakultas
type FakultasListResponse struct {
	Data  []FakultasRefResponse `json:"data"`
	Total int                   `json:"total"`
}

// ProdiListResponse untuk response list prodi
type ProdiListResponse struct {
	Data  []ProdiRefResponse `json:"data"`
	Total int                `json:"total"`
}
