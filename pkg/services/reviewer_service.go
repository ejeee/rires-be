package services

import (
	"errors"
	"fmt"
	"time"

	"rires-be/internal/dto/request"
	"rires-be/internal/dto/response"
	"rires-be/internal/models"
	"rires-be/pkg/database"

	"gorm.io/gorm"
)

// ReviewerService handles reviewer management business logic
type ReviewerService struct {
	externalService *ExternalDataService
}

// NewReviewerService creates a new reviewer service
func NewReviewerService() *ReviewerService {
	return &ReviewerService{
		externalService: NewExternalDataService(),
	}
}

// GetAllReviewers gets all active reviewers from local db_reviewer table
func (s *ReviewerService) GetAllReviewers() ([]response.ReviewerResponse, error) {
	var reviewers []models.Reviewer
	if err := database.DB.Where("hapus = ?", 0).Order("nama_reviewer ASC").Find(&reviewers).Error; err != nil {
		return nil, err
	}

	result := make([]response.ReviewerResponse, 0)
	for _, reviewer := range reviewers {
		result = append(result, response.ReviewerResponse{
			ID:          reviewer.ID,
			IDPegawai:   reviewer.IDPegawai,
			NamaPegawai: reviewer.NamaReviewer, // From local db_reviewer
			NamaLengkap: reviewer.NamaReviewer, // Same value, already has gelar
			EmailUmm:    reviewer.EmailUmm,
			IsActive:    reviewer.IsActive,
			TglInsert:   reviewer.TglInsert,
		})
	}

	return result, nil
}

// GetAvailablePegawai gets all pegawai with email_umm (potential reviewers)
func (s *ReviewerService) GetAvailablePegawai() ([]response.AvailablePegawaiResponse, error) {
	// Get all pegawai with email_umm from SIMPEG
	pegawaiList, err := s.externalService.GetAllReviewers()
	if err != nil {
		return nil, err
	}

	// Get already activated reviewers
	var activatedReviewers []models.Reviewer
	database.DB.Where("hapus = ?", 0).Find(&activatedReviewers)

	// Create map for quick lookup
	activatedMap := make(map[int]bool)
	for _, r := range activatedReviewers {
		activatedMap[r.IDPegawai] = true
	}

	// Build response
	result := make([]response.AvailablePegawaiResponse, 0)
	for _, pegawai := range pegawaiList {
		// Only include pegawai with email_umm
		if pegawai.EmailUMM == "" {
			continue
		}

		result = append(result, response.AvailablePegawaiResponse{
			ID:          pegawai.ID,
			NamaPegawai: pegawai.NamaPegawai,
			NamaLengkap: pegawai.GetNamaLengkap(),
			EmailUmm:    pegawai.EmailUMM,
			IsActivated: activatedMap[pegawai.ID],
		})
	}

	return result, nil
}

// ActivateReviewer activates pegawai as reviewer
func (s *ReviewerService) ActivateReviewer(req *request.ActivateReviewerRequest, userID int) (*response.ReviewerResponse, error) {
	// 1. Check if already activated (hapus = 0)
	var existing models.Reviewer
	err := database.DB.Where("id_pegawai = ? AND hapus = ?", req.IDPegawai, 0).First(&existing).Error
	if err == nil {
		return nil, errors.New("pegawai sudah diaktifkan sebagai reviewer")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2. Get pegawai data from SIMPEG
	pegawai, err := s.externalService.GetPegawaiByID(req.IDPegawai)
	if err != nil {
		return nil, errors.New("pegawai tidak ditemukan di database SIMPEG")
	}

	// 3. Validate email_umm exists
	if pegawai.EmailUMM == "" {
		return nil, errors.New("pegawai tidak memiliki email UMM, tidak dapat diaktifkan sebagai reviewer")
	}

	// 4. Check if soft-deleted record exists (for re-activation)
	var softDeleted models.Reviewer
	userUpdateStr := fmt.Sprintf("%d", userID)

	if err := database.DB.Where("id_pegawai = ? AND hapus = ?", req.IDPegawai, 1).First(&softDeleted).Error; err == nil {
		// Re-activate: update existing soft-deleted record
		updates := map[string]interface{}{
			"nama_reviewer": pegawai.GetNamaLengkap(),
			"email_umm":     pegawai.EmailUMM,
			"is_active":     1,
			"status":        1,
			"hapus":         0,
			"user_update":   userUpdateStr,
		}

		if err := database.DB.Model(&softDeleted).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("gagal mengaktifkan kembali reviewer: %w", err)
		}

		return &response.ReviewerResponse{
			ID:          softDeleted.ID,
			IDPegawai:   softDeleted.IDPegawai,
			NamaPegawai: pegawai.GetNamaLengkap(),
			EmailUmm:    pegawai.EmailUMM,
			IsActive:    1,
			TglInsert:   softDeleted.TglInsert,
		}, nil
	}

	// 5. Create new reviewer record
	now := time.Now()
	reviewer := &models.Reviewer{
		IDPegawai:    pegawai.ID,
		NamaReviewer: pegawai.GetNamaLengkap(), // Full name with gelar
		EmailUmm:     pegawai.EmailUMM,
		IsActive:     1,
		Status:       1,
		Hapus:        0,
		TglInsert:    &now,
		UserUpdate:   userUpdateStr,
	}

	if err := database.DB.Create(reviewer).Error; err != nil {
		return nil, fmt.Errorf("gagal mengaktifkan reviewer: %w", err)
	}

	// 6. Return response
	return &response.ReviewerResponse{
		ID:          reviewer.ID,
		IDPegawai:   reviewer.IDPegawai,
		NamaPegawai: reviewer.NamaReviewer,
		EmailUmm:    reviewer.EmailUmm,
		IsActive:    reviewer.IsActive,
		TglInsert:   reviewer.TglInsert,
	}, nil
}

// UpdateReviewer updates reviewer status
func (s *ReviewerService) UpdateReviewer(id int, req *request.UpdateReviewerRequest, userID int) (*response.ReviewerResponse, error) {
	// 1. Get reviewer
	var reviewer models.Reviewer
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&reviewer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reviewer tidak ditemukan")
		}
		return nil, err
	}

	// 2. Update
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"is_active":   req.IsActive,
		"user_update": userUpdateStr,
	}

	if err := database.DB.Model(&reviewer).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 3. Return response
	return &response.ReviewerResponse{
		ID:          reviewer.ID,
		IDPegawai:   reviewer.IDPegawai,
		NamaPegawai: reviewer.NamaReviewer,
		EmailUmm:    reviewer.EmailUmm,
		IsActive:    req.IsActive,
		TglInsert:   reviewer.TglInsert,
	}, nil
}

// DeleteReviewer soft deletes reviewer
func (s *ReviewerService) DeleteReviewer(id int, userID int) error {
	// 1. Get reviewer
	var reviewer models.Reviewer
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&reviewer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reviewer tidak ditemukan")
		}
		return err
	}

	// 2. Soft delete
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"hapus":       1,
		"user_update": userUpdateStr,
	}

	return database.DB.Model(&reviewer).Updates(updates).Error
}

// IsActiveReviewer checks if pegawai is an active reviewer
func (s *ReviewerService) IsActiveReviewer(idPegawai int) bool {
	var reviewer models.Reviewer
	err := database.DB.Where("id_pegawai = ? AND is_active = ? AND status = ? AND hapus = ?", idPegawai, 1, 1, 0).First(&reviewer).Error
	return err == nil
}
