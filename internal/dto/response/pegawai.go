package response

// PegawaiResponse represents simplified pegawai/reviewer data for API response
type PegawaiResponse struct {
	ID            int              `json:"id"`
	Nama          string           `json:"nama"`
	NamaLengkap   string           `json:"nama_lengkap"` // With gelar
	HP            string           `json:"no_hp"`
	Email         string           `json:"email"`
	EmailUMM	  string           `json:"email_umm"`
	Foto		  string           `json:"foto"`

	Fakultas      *FakultasResponse `json:"fakultas,omitempty"`
}