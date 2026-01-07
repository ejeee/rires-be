package response

import "time"

// ReviewerResponse represents reviewer information
type ReviewerResponse struct {
	ID          int        `json:"id"`
	IDPegawai   int        `json:"id_pegawai"`
	NamaPegawai string     `json:"nama_pegawai"`
	EmailUmm    string     `json:"email_umm"`
	IsActive    int        `json:"is_active"`
	TglInsert   *time.Time `json:"tgl_insert"`
}

// AvailablePegawaiResponse represents pegawai that can be activated as reviewer
type AvailablePegawaiResponse struct {
	ID           int    `json:"id"`
	NamaPegawai  string `json:"nama_pegawai"`
	NamaLengkap  string `json:"nama_lengkap"`
	EmailUmm     string `json:"email_umm"`
	IsActivated  bool   `json:"is_activated"` // Already in db_reviewer?
}