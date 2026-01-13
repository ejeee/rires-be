package utils

import (
	"errors"
	"rires-be/internal/models"
	"rires-be/pkg/database"
	"time"
)

// StatusValidator handles status flow validation
type StatusValidator struct{}

// NewStatusValidator creates a new instance of StatusValidator
func NewStatusValidator() *StatusValidator {
	return &StatusValidator{}
}

// ========================================
// JUDUL VALIDATION
// ========================================

// CanSubmitProposal checks if mahasiswa can submit proposal
// Requirements: status_judul = ACC
func (v *StatusValidator) CanSubmitProposal(pengajuan *models.Pengajuan) error {
	if pengajuan == nil {
		return errors.New("pengajuan is nil")
	}

	if pengajuan.StatusJudul != "ACC" {
		return errors.New("judul belum di-ACC. Status saat ini: " + pengajuan.StatusJudul)
	}

	return nil
}

// CanReviseJudul checks if mahasiswa can revise judul
// Requirements: status_judul = REVISI
func (v *StatusValidator) CanReviseJudul(pengajuan *models.Pengajuan) error {
	if pengajuan == nil {
		return errors.New("pengajuan is nil")
	}

	if pengajuan.StatusJudul != "REVISI" {
		return errors.New("judul tidak dalam status REVISI. Status saat ini: " + pengajuan.StatusJudul)
	}

	return nil
}

// CanReviewJudul checks if reviewer can review judul
// Requirements: status_judul = PENDING or ON_REVIEW
func (v *StatusValidator) CanReviewJudul(pengajuan *models.Pengajuan) error {
	if pengajuan == nil {
		return errors.New("pengajuan is nil")
	}

	validStatuses := []string{"PENDING", "ON_REVIEW", "REVISI"}
	for _, status := range validStatuses {
		if pengajuan.StatusJudul == status {
			return nil
		}
	}

	return errors.New("judul tidak bisa di-review. Status saat ini: " + pengajuan.StatusJudul)
}

// ========================================
// PROPOSAL VALIDATION
// ========================================

// CanReviseProposal checks if mahasiswa can revise proposal
// Requirements: status_proposal = REVISI
func (v *StatusValidator) CanReviseProposal(pengajuan *models.Pengajuan) error {
	if pengajuan == nil {
		return errors.New("pengajuan is nil")
	}

	if pengajuan.StatusProposal != "REVISI" {
		return errors.New("proposal tidak dalam status REVISI. Status saat ini: " + pengajuan.StatusProposal)
	}

	return nil
}

// CanReviewProposal checks if reviewer can review proposal
// Requirements: status_proposal = PENDING or ON_REVIEW, file_proposal exists
func (v *StatusValidator) CanReviewProposal(pengajuan *models.Pengajuan) error {
	if pengajuan == nil {
		return errors.New("pengajuan is nil")
	}

	if pengajuan.FileProposal == "" {
		return errors.New("proposal belum di-upload")
	}

	validStatuses := []string{"PENDING", "ON_REVIEW", "REVISI"}
	for _, status := range validStatuses {
		if pengajuan.StatusProposal == status {
			return nil
		}
	}

	return errors.New("proposal tidak bisa di-review. Status saat ini: " + pengajuan.StatusProposal)
}

// ========================================
// TEAM VALIDATION
// ========================================

// ValidateTeamSize validates team size (max 5: 1 ketua + 4 anggota)
func (v *StatusValidator) ValidateTeamSize(anggota []models.PengajuanAnggota) error {
	if len(anggota) > 5 {
		return errors.New("maksimal tim adalah 5 orang (1 ketua + 4 anggota)")
	}

	if len(anggota) == 0 {
		return errors.New("tim harus memiliki minimal 1 anggota (ketua)")
	}

	return nil
}

// ValidateTeamStructure validates team structure (must have 1 ketua)
func (v *StatusValidator) ValidateTeamStructure(anggota []models.PengajuanAnggota) error {
	ketuaCount := 0
	for _, member := range anggota {
		if member.IsKetua == 1 {
			ketuaCount++
		}
	}

	if ketuaCount == 0 {
		return errors.New("tim harus memiliki 1 ketua")
	}

	if ketuaCount > 1 {
		return errors.New("tim hanya boleh memiliki 1 ketua")
	}

	return nil
}

// ValidateNoDuplicateNIM validates no duplicate NIM in team
func (v *StatusValidator) ValidateNoDuplicateNIM(anggota []models.PengajuanAnggota) error {
	nimMap := make(map[string]bool)

	for _, member := range anggota {
		if nimMap[member.NIM] {
			return errors.New("NIM duplikat ditemukan: " + member.NIM)
		}
		nimMap[member.NIM] = true
	}

	return nil
}

// ========================================
// REGISTRATION PERIOD VALIDATION
// ========================================

// CanSubmitPengajuan checks if registration period is open
func (v *StatusValidator) CanSubmitPengajuan() error {
	// Check if tgl_setting is active
	var setting models.TglSetting
	err := database.DB.Where("is_active = ? AND status = ? AND hapus = ?", 1, 1, 0).
		First(&setting).Error

	if err != nil {
		return errors.New("pendaftaran PKM sedang ditutup")
	}

	now := time.Now()

	// Check if current date is within the registration period
	if now.Before(setting.TglDaftarAwal) {
		return errors.New("pendaftaran belum dibuka. Akan dibuka pada: " + setting.TglDaftarAwal.Format("02 January 2006"))
	}

	if now.After(setting.TglDaftarAkhir) {
		return errors.New("pendaftaran sudah ditutup pada: " + setting.TglDaftarAkhir.Format("02 January 2006"))
	}

	return nil
}

// ========================================
// OWNERSHIP VALIDATION
// ========================================

// IsOwner checks if given NIM is the owner (ketua) of pengajuan
func (v *StatusValidator) IsOwner(pengajuan *models.Pengajuan, nim string) bool {
	if pengajuan == nil {
		return false
	}
	return pengajuan.NIMKetua == nim
}

// IsMemberOfTeam checks if given NIM is member of the team
func (v *StatusValidator) IsMemberOfTeam(anggota []models.PengajuanAnggota, nim string) bool {
	for _, member := range anggota {
		if member.NIM == nim {
			return true
		}
	}
	return false
}