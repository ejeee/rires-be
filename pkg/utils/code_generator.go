package utils

import (
	"fmt"
	"rires-be/internal/models"
	"rires-be/pkg/database"
	"time"
)

// GenerateKodePengajuan generates unique code for pengajuan
// Format: PKM-{KODE_KATEGORI}-{TAHUN}-{SEQUENCE}
// Example: PKM-K-2026-001, PKM-RE-2026-002
func GenerateKodePengajuan(kategori *models.KategoriPKM, tahun int) (string, error) {
	if kategori == nil {
		return "", fmt.Errorf("kategori is required")
	}

	if tahun == 0 {
		tahun = time.Now().Year()
	}

	// Extract kategori code (e.g., "PKM-K" -> "K", "PKM-RE" -> "RE")
	kodeKategori := extractKodeKategori(kategori.NamaKategori)

	// Get last sequence number for this kategori and tahun
	sequence, err := getLastSequence(kategori.ID, tahun)
	if err != nil {
		return "", err
	}

	// Increment sequence
	sequence++

	// Generate code: PKM-K-2026-001
	code := fmt.Sprintf("PKM-%s-%d-%03d", kodeKategori, tahun, sequence)

	return code, nil
}

// extractKodeKategori extracts the kategori code from kategori name
// PKM-K -> K
// PKM-RE -> RE
// PKM-RSH -> RSH
func extractKodeKategori(namaKategori string) string {
	// Remove "PKM-" prefix
	if len(namaKategori) > 4 && namaKategori[:4] == "PKM-" {
		return namaKategori[4:]
	}
	return namaKategori
}

// getLastSequence gets the last sequence number for given kategori and tahun
func getLastSequence(kategoriID, tahun int) (int, error) {
	var lastPengajuan models.Pengajuan

	err := database.DB.Where("id_kategori = ? AND tahun = ? AND hapus = ?", kategoriID, tahun, 0).
		Order("id DESC").
		First(&lastPengajuan).Error

	if err != nil {
		// No records found, start from 0
		return 0, nil
	}

	// Extract sequence from kode_pengajuan (last 3 digits)
	// PKM-K-2026-001 -> 001
	var sequence int
	fmt.Sscanf(lastPengajuan.KodePengajuan[len(lastPengajuan.KodePengajuan)-3:], "%d", &sequence)

	return sequence, nil
}

// ValidateKodePengajuan checks if kode_pengajuan is unique
func ValidateKodePengajuan(kodePengajuan string) bool {
	var count int64
	database.DB.Model(&models.Pengajuan{}).
		Where("kode_pengajuan = ? AND hapus = ?", kodePengajuan, 0).
		Count(&count)

	return count == 0 // True if unique (count = 0)
}