package response

// LoginResponse adalah struktur untuk response login
type LoginResponse struct {
	Token     string      `json:"token"`
	UserType  string      `json:"user_type"` // admin, mahasiswa, pegawai
	User      interface{} `json:"user"`
	ExpiresIn int         `json:"expires_in"` // dalam jam
}

// AdminLoginResponse data admin setelah login
type AdminLoginResponse struct {
	ID        int    `json:"id"`
	NamaUser  string `json:"nama_user"`
	Username  string `json:"username"`
	LevelUser int    `json:"level_user"`
	Status    int    `json:"status"`
}

// MahasiswaLoginResponse data mahasiswa dari API external
type MahasiswaLoginResponse struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Prodi    string `json:"prodi"`
	Fakultas string `json:"fakultas"`
	Email    string `json:"email"`
	// Tambahkan field lain sesuai response API
}

// PegawaiLoginResponse data pegawai dari API external
type PegawaiLoginResponse struct {
	NIP     string `json:"nip"`
	Nama    string `json:"nama"`
	Jabatan string `json:"jabatan"`
	Unit    string `json:"unit"`
	Email   string `json:"email"`
	// Tambahkan field lain sesuai response API
}