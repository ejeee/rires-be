package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword meng-hash password menggunakan bcrypt
func HashPassword(password string) (string, error) {
	// Generate hash dengan cost 14 (recommended untuk security)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword memverifikasi password dengan hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}