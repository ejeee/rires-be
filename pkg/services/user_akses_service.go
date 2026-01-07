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

// UserAksesService handles user access management business logic
type UserAksesService struct{}

// NewUserAksesService creates a new user akses service
func NewUserAksesService() *UserAksesService {
	return &UserAksesService{}
}

// GetAllAccesses gets all user accesses with optional filters
func (s *UserAksesService) GetAllAccesses(idUserLevel int, idMenu int) ([]response.UserAksesResponse, error) {
	query := database.DB.Where("hapus = ?", 0)

	// Apply filters
	if idUserLevel > 0 {
		query = query.Where("id_user_level = ?", idUserLevel)
	}
	if idMenu > 0 {
		query = query.Where("id_menu = ?", idMenu)
	}

	var accesses []models.UserAkses
	if err := query.Preload("UserLevel").Preload("Menu").Order("id_user_level ASC, id_menu ASC").Find(&accesses).Error; err != nil {
		return nil, err
	}

	// Map to response
	result := make([]response.UserAksesResponse, 0)
	for _, access := range accesses {
		resp := response.UserAksesResponse{
			ID:          access.ID,
			IDUserLevel: access.IDUserLevel,
			IDMenu:      access.IDMenu,
			CanCreate:   access.CanCreate,
			CanUpdate:   access.CanUpdate,
			CanDelete:   access.CanDelete,
			Status:      access.Status,
			TglInsert:   access.TglInsert,
		}

		// Map UserLevel
		if access.UserLevel != nil {
			resp.UserLevel = &response.UserLevelResponse{
				ID:        access.UserLevel.ID,
				NamaLevel: access.UserLevel.NamaLevel,
			}
		}

		// Map Menu
		if access.Menu != nil {
			resp.Menu = &response.MenuSimpleResponse{
				ID:       access.Menu.ID,
				NamaMenu: access.Menu.NamaMenu,
				URLMenu:  access.Menu.URLMenu,
				Lucide:   access.Menu.Lucide,
			}
		}

		result = append(result, resp)
	}

	return result, nil
}

// GetAccessesByUserLevel gets all accesses for a specific user level
func (s *UserAksesService) GetAccessesByUserLevel(idUserLevel int) (*response.UserAksesGroupedResponse, error) {
	// Get user level
	var userLevel models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", idUserLevel, 0).First(&userLevel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user level tidak ditemukan")
		}
		return nil, err
	}

	// Get accesses
	accesses, err := s.GetAllAccesses(idUserLevel, 0)
	if err != nil {
		return nil, err
	}

	return &response.UserAksesGroupedResponse{
		IDUserLevel: userLevel.ID,
		NamaLevel:   userLevel.NamaLevel,
		Accesses:    accesses,
	}, nil
}

// GetAccessDetail gets single access detail
func (s *UserAksesService) GetAccessDetail(id int) (*response.UserAksesResponse, error) {
	var access models.UserAkses
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).
		Preload("UserLevel").
		Preload("Menu").
		First(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akses tidak ditemukan")
		}
		return nil, err
	}

	resp := &response.UserAksesResponse{
		ID:          access.ID,
		IDUserLevel: access.IDUserLevel,
		IDMenu:      access.IDMenu,
		CanCreate:   access.CanCreate,
		CanUpdate:   access.CanUpdate,
		CanDelete:   access.CanDelete,
		Status:      access.Status,
		TglInsert:   access.TglInsert,
	}

	// Map UserLevel
	if access.UserLevel != nil {
		resp.UserLevel = &response.UserLevelResponse{
			ID:        access.UserLevel.ID,
			NamaLevel: access.UserLevel.NamaLevel,
		}
	}

	// Map Menu
	if access.Menu != nil {
		resp.Menu = &response.MenuSimpleResponse{
			ID:       access.Menu.ID,
			NamaMenu: access.Menu.NamaMenu,
			URLMenu:  access.Menu.URLMenu,
			Lucide:   access.Menu.Lucide,
		}
	}

	return resp, nil
}

// CreateAccess creates new user access
func (s *UserAksesService) CreateAccess(req *request.CreateUserAksesRequest, userID int) (*response.UserAksesResponse, error) {
	// 1. Validate user level exists
	var userLevel models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", req.IDUserLevel, 0).First(&userLevel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user level tidak ditemukan")
		}
		return nil, err
	}

	// 2. Validate menu exists
	var menu models.Menu
	if err := database.DB.Where("id = ? AND hapus = ?", req.IDMenu, 0).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu tidak ditemukan")
		}
		return nil, err
	}

	// 3. Check if access already exists
	var existing models.UserAkses
	err := database.DB.Where("id_user_level = ? AND id_menu = ? AND hapus = ?", req.IDUserLevel, req.IDMenu, 0).First(&existing).Error
	if err == nil {
		return nil, errors.New("akses untuk user level dan menu ini sudah ada")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 4. Create access
	now := time.Now()
	userUpdateStr := fmt.Sprintf("%d", userID)

	access := &models.UserAkses{
		IDUserLevel: req.IDUserLevel,
		IDMenu:      req.IDMenu,
		CanCreate:   req.CanCreate,
		CanUpdate:   req.CanUpdate,
		CanDelete:   req.CanDelete,
		Status:      1,
		Hapus:       0,
		TglInsert:   &now,
		UserUpdate:  userUpdateStr,
	}

	if err := database.DB.Create(access).Error; err != nil {
		return nil, fmt.Errorf("gagal membuat akses: %w", err)
	}

	// 5. Return detail
	return s.GetAccessDetail(access.ID)
}

// UpdateAccess updates existing access
func (s *UserAksesService) UpdateAccess(id int, req *request.UpdateUserAksesRequest, userID int) (*response.UserAksesResponse, error) {
	// 1. Get access
	var access models.UserAkses
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akses tidak ditemukan")
		}
		return nil, err
	}

	// 2. Update
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"can_create":  req.CanCreate,
		"can_update":  req.CanUpdate,
		"can_delete":  req.CanDelete,
		"user_update": userUpdateStr,
	}

	if err := database.DB.Model(&access).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 3. Return detail
	return s.GetAccessDetail(id)
}

// DeleteAccess soft deletes access
func (s *UserAksesService) DeleteAccess(id int, userID int) error {
	// 1. Get access
	var access models.UserAkses
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("akses tidak ditemukan")
		}
		return err
	}

	// 2. Soft delete
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"hapus":       1,
		"user_update": userUpdateStr,
	}

	return database.DB.Model(&access).Updates(updates).Error
}

// BulkCreateAccess creates multiple accesses at once for a user level
func (s *UserAksesService) BulkCreateAccess(req *request.BulkCreateUserAksesRequest, userID int) ([]response.UserAksesResponse, error) {
	// 1. Validate user level exists
	var userLevel models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", req.IDUserLevel, 0).First(&userLevel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user level tidak ditemukan")
		}
		return nil, err
	}

	// 2. START TRANSACTION
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	userUpdateStr := fmt.Sprintf("%d", userID)
	createdIDs := make([]int, 0)

	// 3. Create each access
	for _, menuPerm := range req.Menus {
		// Validate menu exists
		var menu models.Menu
		if err := tx.Where("id = ? AND hapus = ?", menuPerm.IDMenu, 0).First(&menu).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("menu dengan id %d tidak ditemukan", menuPerm.IDMenu)
			}
			return nil, err
		}

		// Check if already exists
		var existing models.UserAkses
		err := tx.Where("id_user_level = ? AND id_menu = ? AND hapus = ?", req.IDUserLevel, menuPerm.IDMenu, 0).First(&existing).Error
		if err == nil {
			// Skip if already exists
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, err
		}

		// Create access
		access := &models.UserAkses{
			IDUserLevel: req.IDUserLevel,
			IDMenu:      menuPerm.IDMenu,
			CanCreate:   menuPerm.CanCreate,
			CanUpdate:   menuPerm.CanUpdate,
			CanDelete:   menuPerm.CanDelete,
			Status:      1,
			Hapus:       0,
			TglInsert:   &now,
			UserUpdate:  userUpdateStr,
		}

		if err := tx.Create(access).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal membuat akses untuk menu %d: %w", menuPerm.IDMenu, err)
		}

		createdIDs = append(createdIDs, access.ID)
	}

	// 4. COMMIT
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 5. Return all created accesses
	result := make([]response.UserAksesResponse, 0)
	for _, id := range createdIDs {
		detail, _ := s.GetAccessDetail(id)
		if detail != nil {
			result = append(result, *detail)
		}
	}

	return result, nil
}