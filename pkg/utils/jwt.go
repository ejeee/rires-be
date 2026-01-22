package utils

import (
	"errors"
	"rires-be/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims adalah struktur claims untuk JWT
type JWTClaims struct {
	UserID   uint              `json:"id_user"`
	Email    string            `json:"email"`
	Username string            `json:"username"`
	UserType string            `json:"user_type"` // admin, mahasiswa, pegawai
	UserData map[string]string `json:"user_data"` // Additional user data
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT token baru
func GenerateToken(userID uint, email string) (string, error) {
	// Parse JWT expired hours dari config
	expiredHours, err := strconv.Atoi(config.AppConfig.JWTExpiredHours)
	if err != nil {
		expiredHours = 24 // Default 24 jam
	}

	// Buat claims
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateTokenWithClaims membuat JWT token dengan custom claims
func GenerateTokenWithClaims(userID uint, username, email, userType string, userData map[string]string) (string, error) {
	// Parse JWT expired hours dari config
	expiredHours, err := strconv.Atoi(config.AppConfig.JWTExpiredHours)
	if err != nil {
		expiredHours = 24 // Default 24 jam
	}

	// Buat claims
	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		UserType: userType,
		UserData: userData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi JWT token dan mengembalikan claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Ambil claims dari token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
