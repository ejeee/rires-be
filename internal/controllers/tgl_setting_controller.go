package controllers

import (
	"math"
	"rires-be/internal/dto/request"
	"rires-be/internal/dto/response"
	"rires-be/internal/models"
	"rires-be/pkg/database"
	"rires-be/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TglSettingController struct{}

func NewTglSettingController() *TglSettingController {
	return &TglSettingController{}
}

// GetActive godoc
// @Summary Get Active Registration Period
// @Description Get currently active registration period settings
// @Tags Tanggal Setting
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string,data=response.RegistrationStatusResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tgl-setting/active [get]
func (ctrl *TglSettingController) GetActive(c *fiber.Ctx) error {
	var setting models.TglSetting
	err := database.DB.Where("is_active = ? AND status = ? AND hapus = ?", 1, 1, 0).
		First(&setting).Error

	if err != nil {
		return utils.SuccessResponse(c, "Registration is currently closed", response.RegistrationStatusResponse{
			IsOpen:  false,
			Message: "Pendaftaran PKM sedang ditutup. Silakan cek kembali nanti.",
		})
	}

	now := time.Now()
	isOpen := setting.IsRegistrationOpen()

	var message string
	var daysRemaining int

	if isOpen {
		duration := setting.TglDaftarAkhir.Sub(now)
		daysRemaining = int(duration.Hours() / 24)
		message = "Pendaftaran PKM sedang dibuka!"
	} else if now.Before(setting.TglDaftarAwal) {
		message = "Pendaftaran belum dibuka. Akan dibuka pada " + setting.TglDaftarAwal.Format("02 January 2006")
	} else {
		message = "Pendaftaran sudah ditutup."
	}

	result := response.RegistrationStatusResponse{
		IsOpen:         isOpen,
		Message:        message,
		TglDaftarAwal:  setting.TglDaftarAwal,
		TglDaftarAkhir: setting.TglDaftarAkhir,
		TglReviewAwal:  setting.TglReviewAwal,
		TglReviewAkhir: setting.TglReviewAkhir,
		TglPengumuman:  setting.TglPengumuman,
		DaysRemaining:  daysRemaining,
	}

	return utils.SuccessResponse(c, message, result)
}

// GetList godoc
// @Summary List Tanggal Setting
// @Description Get list of registration period settings with pagination (admin only)
// @Tags Tanggal Setting
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} object{success=bool,message=string,data=response.TglSettingListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tgl-setting [get]
func (ctrl *TglSettingController) GetList(c *fiber.Ctx) error {
	// Pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Query builder
	query := database.DB.Model(&models.TglSetting{}).Where("hapus = ?", 0)

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data
	var settings []models.TglSetting
	if err := query.Order("id DESC").Offset(offset).Limit(perPage).Find(&settings).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Transform to response
	now := time.Now()
	var data []response.TglSettingResponse
	for _, setting := range settings {
		statusText := "Aktif"
		if setting.Status == 2 {
			statusText = "Tidak Aktif"
		}

		isActiveText := "Tidak Aktif"
		if setting.IsActive == 1 {
			isActiveText = "Aktif"
		}

		isRegOpen := setting.IsRegistrationOpen()
		isReviewPeriod := setting.IsReviewPeriod()
		isAnnounced := setting.IsAfterAnnouncement()

		var daysRemaining int
		if isRegOpen {
			duration := setting.TglDaftarAkhir.Sub(now)
			daysRemaining = int(duration.Hours() / 24)
		}

		data = append(data, response.TglSettingResponse{
			ID:             setting.ID,
			TglDaftarAwal:  setting.TglDaftarAwal,
			TglDaftarAkhir: setting.TglDaftarAkhir,
			TglReviewAwal:  setting.TglReviewAwal,
			TglReviewAkhir: setting.TglReviewAkhir,
			TglPengumuman:  setting.TglPengumuman,
			Keterangan:     setting.Keterangan,
			IsActive:       setting.IsActive,
			IsActiveText:   isActiveText,
			IsRegOpen:      isRegOpen,
			IsReviewPeriod: isReviewPeriod,
			IsAnnounced:    isAnnounced,
			Status:         setting.Status,
			StatusText:     statusText,
			DaysRemaining:  daysRemaining,
			TglInsert:      setting.TglInsert,
			TglUpdate:      setting.TglUpdate,
			UserUpdate:     setting.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.TglSettingListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetByID godoc
// @Summary Get Tanggal Setting by ID
// @Description Get registration period setting detail by ID
// @Tags Tanggal Setting
// @Accept json
// @Produce json
// @Param id path int true "Tanggal Setting ID"
// @Success 200 {object} object{success=bool,message=string,data=response.TglSettingResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tgl-setting/{id} [get]
func (ctrl *TglSettingController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var setting models.TglSetting
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&setting).Error; err != nil {
		return utils.NotFoundResponse(c, "Tanggal setting not found")
	}

	statusText := "Aktif"
	if setting.Status == 2 {
		statusText = "Tidak Aktif"
	}

	isActiveText := "Tidak Aktif"
	if setting.IsActive == 1 {
		isActiveText = "Aktif"
	}

	isRegOpen := setting.IsRegistrationOpen()
	isReviewPeriod := setting.IsReviewPeriod()
	isAnnounced := setting.IsAfterAnnouncement()

	now := time.Now()
	var daysRemaining int
	if isRegOpen {
		duration := setting.TglDaftarAkhir.Sub(now)
		daysRemaining = int(duration.Hours() / 24)
	}

	result := response.TglSettingResponse{
		ID:             setting.ID,
		TglDaftarAwal:  setting.TglDaftarAwal,
		TglDaftarAkhir: setting.TglDaftarAkhir,
		TglReviewAwal:  setting.TglReviewAwal,
		TglReviewAkhir: setting.TglReviewAkhir,
		TglPengumuman:  setting.TglPengumuman,
		Keterangan:     setting.Keterangan,
		IsActive:       setting.IsActive,
		IsActiveText:   isActiveText,
		IsRegOpen:      isRegOpen,
		IsReviewPeriod: isReviewPeriod,
		IsAnnounced:    isAnnounced,
		Status:         setting.Status,
		StatusText:     statusText,
		DaysRemaining:  daysRemaining,
		TglInsert:      setting.TglInsert,
		TglUpdate:      setting.TglUpdate,
		UserUpdate:     setting.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create Tanggal Setting
// @Description Create new registration period setting and set as active (deactivate others)
// @Tags Tanggal Setting
// @Accept json
// @Produce json
// @Param body body request.CreateTglSettingRequest true "Tanggal Setting Data"
// @Success 201 {object} object{success=bool,message=string,data=response.TglSettingResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tgl-setting [post]
func (ctrl *TglSettingController) Create(c *fiber.Ctx) error {
	var req request.CreateTglSettingRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.TglDaftarAwal == "" || req.TglDaftarAkhir == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Parse dates (format: YYYY-MM-DD)
	tglDaftarAwal, err := time.Parse("2006-01-02", req.TglDaftarAwal)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tgl_daftar_awal format. Use: YYYY-MM-DD")
	}

	tglDaftarAkhir, err := time.Parse("2006-01-02", req.TglDaftarAkhir)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tgl_daftar_akhir format. Use: YYYY-MM-DD")
	}

	// Validate date range
	if tglDaftarAkhir.Before(tglDaftarAwal) {
		return utils.BadRequestResponse(c, "Tgl daftar akhir must be after tgl daftar awal")
	}

	// Parse optional dates
	var tglReviewAwal, tglReviewAkhir, tglPengumuman time.Time

	if req.TglReviewAwal != "" {
		tglReviewAwal, err = time.Parse("2006-01-02", req.TglReviewAwal)
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid tgl_review_awal format. Use: YYYY-MM-DD")
		}
	}

	if req.TglReviewAkhir != "" {
		tglReviewAkhir, err = time.Parse("2006-01-02", req.TglReviewAkhir)
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid tgl_review_akhir format. Use: YYYY-MM-DD")
		}
	}

	if req.TglPengumuman != "" {
		tglPengumuman, err = time.Parse("2006-01-02", req.TglPengumuman)
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid tgl_pengumuman format. Use: YYYY-MM-DD")
		}
	}

	// Validate review dates if both are provided
	if !tglReviewAwal.IsZero() && !tglReviewAkhir.IsZero() {
		if tglReviewAkhir.Before(tglReviewAwal) {
			return utils.BadRequestResponse(c, "Tgl review akhir must be after tgl review awal")
		}
	}

	// Deactivate all other periods
	database.DB.Model(&models.TglSetting{}).
		Where("is_active = ?", 1).
		Update("is_active", 0)

	// Create new setting (automatically set as active)
	now := time.Now()
	setting := models.TglSetting{
		TglDaftarAwal:  tglDaftarAwal,
		TglDaftarAkhir: tglDaftarAkhir,
		TglReviewAwal:  tglReviewAwal,
		TglReviewAkhir: tglReviewAkhir,
		TglPengumuman:  tglPengumuman,
		Keterangan:     req.Keterangan,
		IsActive:       1, // Always set new setting as active
		Status:         req.Status,
		Hapus:          0,
		TglInsert:      &now,
		UserUpdate:     strconv.Itoa(int(utils.GetCurrentUserID(c))),
	}

	if err := database.DB.Create(&setting).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create tanggal setting")
	}

	statusText := "Aktif"
	if setting.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.TglSettingResponse{
		ID:             setting.ID,
		TglDaftarAwal:  setting.TglDaftarAwal,
		TglDaftarAkhir: setting.TglDaftarAkhir,
		TglReviewAwal:  setting.TglReviewAwal,
		TglReviewAkhir: setting.TglReviewAkhir,
		TglPengumuman:  setting.TglPengumuman,
		Keterangan:     setting.Keterangan,
		IsActive:       setting.IsActive,
		IsActiveText:   "Aktif",
		IsRegOpen:      setting.IsRegistrationOpen(),
		IsReviewPeriod: setting.IsReviewPeriod(),
		IsAnnounced:    setting.IsAfterAnnouncement(),
		Status:         setting.Status,
		StatusText:     statusText,
		TglInsert:      setting.TglInsert,
		TglUpdate:      setting.TglUpdate,
		UserUpdate:     setting.UserUpdate,
	}

	return utils.CreatedResponse(c, "Tanggal setting created and set as active", result)
}

// Update godoc
// @Summary Update Tanggal Setting
// @Description Update existing registration period setting
// @Tags Tanggal Setting
// @Accept json
// @Produce json
// @Param id path int true "Tanggal Setting ID"
// @Param body body request.UpdateTglSettingRequest true "Tanggal Setting Data"
// @Success 200 {object} object{success=bool,message=string,data=response.TglSettingResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tgl-setting/{id} [put]
func (ctrl *TglSettingController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateTglSettingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.TglDaftarAwal == "" || req.TglDaftarAkhir == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Find existing
	var setting models.TglSetting
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&setting).Error; err != nil {
		return utils.NotFoundResponse(c, "Tanggal setting not found")
	}

	// Parse dates
	tglDaftarAwal, err := time.Parse("2006-01-02", req.TglDaftarAwal)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tgl_daftar_awal format. Use: YYYY-MM-DD")
	}

	tglDaftarAkhir, err := time.Parse("2006-01-02", req.TglDaftarAkhir)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tgl_daftar_akhir format. Use: YYYY-MM-DD")
	}

	// Validate date range
	if tglDaftarAkhir.Before(tglDaftarAwal) {
		return utils.BadRequestResponse(c, "Tgl daftar akhir must be after tgl daftar awal")
	}

	// Parse optional dates
	var tglReviewAwal, tglReviewAkhir, tglPengumuman time.Time

	if req.TglReviewAwal != "" {
		tglReviewAwal, err = time.Parse("2006-01-02", req.TglReviewAwal)
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid tgl_review_awal format. Use: YYYY-MM-DD")
		}
	}

	if req.TglReviewAkhir != "" {
		tglReviewAkhir, err = time.Parse("2006-01-02", req.TglReviewAkhir)
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid tgl_review_akhir format. Use: YYYY-MM-DD")
		}
	}

	if req.TglPengumuman != "" {
		tglPengumuman, err = time.Parse("2006-01-02", req.TglPengumuman)
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid tgl_pengumuman format. Use: YYYY-MM-DD")
		}
	}

	// Validate review dates if both are provided
	if !tglReviewAwal.IsZero() && !tglReviewAkhir.IsZero() {
		if tglReviewAkhir.Before(tglReviewAwal) {
			return utils.BadRequestResponse(c, "Tgl review akhir must be after tgl review awal")
		}
	}

	// Update
	setting.TglDaftarAwal = tglDaftarAwal
	setting.TglDaftarAkhir = tglDaftarAkhir
	setting.TglReviewAwal = tglReviewAwal
	setting.TglReviewAkhir = tglReviewAkhir
	setting.TglPengumuman = tglPengumuman
	setting.Keterangan = req.Keterangan
	setting.Status = req.Status
	setting.UserUpdate = strconv.Itoa(int(utils.GetCurrentUserID(c)))

	if err := database.DB.Save(&setting).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update tanggal setting")
	}

	statusText := "Aktif"
	if setting.Status == 2 {
		statusText = "Tidak Aktif"
	}

	isActiveText := "Tidak Aktif"
	if setting.IsActive == 1 {
		isActiveText = "Aktif"
	}

	result := response.TglSettingResponse{
		ID:             setting.ID,
		TglDaftarAwal:  setting.TglDaftarAwal,
		TglDaftarAkhir: setting.TglDaftarAkhir,
		TglReviewAwal:  setting.TglReviewAwal,
		TglReviewAkhir: setting.TglReviewAkhir,
		TglPengumuman:  setting.TglPengumuman,
		Keterangan:     setting.Keterangan,
		IsActive:       setting.IsActive,
		IsActiveText:   isActiveText,
		IsRegOpen:      setting.IsRegistrationOpen(),
		IsReviewPeriod: setting.IsReviewPeriod(),
		IsAnnounced:    setting.IsAfterAnnouncement(),
		Status:         setting.Status,
		StatusText:     statusText,
		TglInsert:      setting.TglInsert,
		TglUpdate:      setting.TglUpdate,
		UserUpdate:     setting.UserUpdate,
	}

	return utils.SuccessResponse(c, "Tanggal setting updated successfully", result)
}

// Delete godoc
// @Summary Delete Tanggal Setting
// @Description Soft delete registration period setting
// @Tags Tanggal Setting
// @Accept json
// @Produce json
// @Param id path int true "Tanggal Setting ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tgl-setting/{id} [delete]
func (ctrl *TglSettingController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var setting models.TglSetting
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&setting).Error; err != nil {
		return utils.NotFoundResponse(c, "Tanggal setting not found")
	}

	// Soft delete
	setting.Hapus = 1
	setting.IsActive = 0 // Also deactivate
	setting.UserUpdate = strconv.Itoa(int(utils.GetCurrentUserID(c)))

	if err := database.DB.Save(&setting).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete tanggal setting")
	}

	return utils.SuccessResponse(c, "Tanggal setting deleted successfully", nil)
}
