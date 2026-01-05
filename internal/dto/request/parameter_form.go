package request

// CreateParameterFormRequest untuk create parameter form
type CreateParameterFormRequest struct {
	KategoriID    int    `json:"kategori_id" validate:"required"`
	NamaParameter string `json:"nama_parameter" validate:"required,min=2,max=100"`
	Label         string `json:"label" validate:"required,min=3,max=200"`
	TipeInput     string `json:"tipe_input" validate:"required"` // text, textarea, number, file, radio, select, etc
	Validasi      string `json:"validasi"`                       // JSON string
	Placeholder   string `json:"placeholder"`
	HelpText      string `json:"help_text"`
	Opsi          string `json:"opsi"` // JSON string for radio/select options
	Urutan        int    `json:"urutan" validate:"min=0"`
	Status        int    `json:"status" validate:"required,oneof=1 2"` // 1=active, 2=inactive
}

// UpdateParameterFormRequest untuk update parameter form
type UpdateParameterFormRequest struct {
	KategoriID    int    `json:"kategori_id" validate:"required"`
	NamaParameter string `json:"nama_parameter" validate:"required,min=2,max=100"`
	Label         string `json:"label" validate:"required,min=3,max=200"`
	TipeInput     string `json:"tipe_input" validate:"required"`
	Validasi      string `json:"validasi"`
	Placeholder   string `json:"placeholder"`
	HelpText      string `json:"help_text"`
	Opsi          string `json:"opsi"` // JSON string for radio/select options
	Urutan        int    `json:"urutan" validate:"min=0"`
	Status        int    `json:"status" validate:"required,oneof=1 2"`
}