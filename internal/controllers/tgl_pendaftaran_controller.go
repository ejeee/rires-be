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

type TanggalPendaftaranController struct{}

func NewTanggalPendaftaranController() *TanggalPendaftaranController {
	return &TanggalPendaftaranController{}
}

// GetActive godoc
// @Summary Get Active Registration Period
// @Description Get currently active registration period (for checking if registration is open)
// @Tags Tanggal Pendaftaran
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string,data=response.RegistrationStatusResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tanggal-pendaftaran/active [get]
func (ctrl *TanggalPendaftaranController) GetActive(c *fiber.Ctx) error {
	var period models.TanggalPendaftaran
	err := database.DB.Where("is_active = ? AND status = ? AND hapus = ?", 1, 1, 0).
		First(&period).Error

	if err != nil {
		return utils.SuccessResponse(c, "Registration is currently closed", response.RegistrationStatusResponse{
			IsOpen:  false,
			Message: "Pendaftaran PKM sedang ditutup. Silakan cek kembali nanti.",
		})
	}

	now := time.Now()
	isOpen := period.IsOpen()
	
	var message string
	var daysRemaining int
	
	if isOpen {
		duration := period.TanggalSelesai.Sub(now)
		daysRemaining = int(duration.Hours() / 24)
		message = "Pendaftaran PKM sedang dibuka!"
	} else if now.Before(period.TanggalMulai) {
		message = "Pendaftaran belum dibuka. Akan dibuka pada " + period.TanggalMulai.Format("02 January 2006")
	} else {
		message = "Pendaftaran sudah ditutup."
	}

	result := response.RegistrationStatusResponse{
		IsOpen:         isOpen,
		Message:        message,
		TanggalMulai:   period.TanggalMulai,
		TanggalSelesai: period.TanggalSelesai,
		DaysRemaining:  daysRemaining,
	}

	return utils.SuccessResponse(c, message, result)
}

// GetList godoc
// @Summary List Tanggal Pendaftaran
// @Description Get list of registration periods with pagination (admin only)
// @Tags Tanggal Pendaftaran
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} object{success=bool,message=string,data=response.TanggalPendaftaranListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tanggal-pendaftaran [get]
func (ctrl *TanggalPendaftaranController) GetList(c *fiber.Ctx) error {
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
	query := database.DB.Model(&models.TanggalPendaftaran{}).Where("hapus = ?", 0)

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data
	var periods []models.TanggalPendaftaran
	if err := query.Order("id DESC").Offset(offset).Limit(perPage).Find(&periods).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Transform to response
	now := time.Now()
	var data []response.TanggalPendaftaranResponse
	for _, period := range periods {
		statusText := "Aktif"
		if period.Status == 2 {
			statusText = "Tidak Aktif"
		}

		isActiveText := "Tidak Aktif"
		if period.IsActive == 1 {
			isActiveText = "Aktif"
		}

		isOpen := period.IsOpen()
		
		var daysRemaining int
		if isOpen {
			duration := period.TanggalSelesai.Sub(now)
			daysRemaining = int(duration.Hours() / 24)
		}

		data = append(data, response.TanggalPendaftaranResponse{
			ID:              period.ID,
			TanggalMulai:    period.TanggalMulai,
			TanggalSelesai:  period.TanggalSelesai,
			Keterangan:      period.Keterangan,
			IsActive:        period.IsActive,
			IsActiveText:    isActiveText,
			IsOpen:          isOpen,
			Status:          period.Status,
			StatusText:      statusText,
			DaysRemaining:   daysRemaining,
			TglInsert:       period.TglInsert,
			TglUpdate:       period.TglUpdate,
			UserUpdate:      period.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.TanggalPendaftaranListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetByID godoc
// @Summary Get Tanggal Pendaftaran by ID
// @Description Get registration period detail by ID
// @Tags Tanggal Pendaftaran
// @Accept json
// @Produce json
// @Param id path int true "Tanggal Pendaftaran ID"
// @Success 200 {object} object{success=bool,message=string,data=response.TanggalPendaftaranResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tanggal-pendaftaran/{id} [get]
func (ctrl *TanggalPendaftaranController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var period models.TanggalPendaftaran
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&period).Error; err != nil {
		return utils.NotFoundResponse(c, "Tanggal pendaftaran not found")
	}

	statusText := "Aktif"
	if period.Status == 2 {
		statusText = "Tidak Aktif"
	}

	isActiveText := "Tidak Aktif"
	if period.IsActive == 1 {
		isActiveText = "Aktif"
	}

	isOpen := period.IsOpen()
	
	now := time.Now()
	var daysRemaining int
	if isOpen {
		duration := period.TanggalSelesai.Sub(now)
		daysRemaining = int(duration.Hours() / 24)
	}

	result := response.TanggalPendaftaranResponse{
		ID:              period.ID,
		TanggalMulai:    period.TanggalMulai,
		TanggalSelesai:  period.TanggalSelesai,
		Keterangan:      period.Keterangan,
		IsActive:        period.IsActive,
		IsActiveText:    isActiveText,
		IsOpen:          isOpen,
		Status:          period.Status,
		StatusText:      statusText,
		DaysRemaining:   daysRemaining,
		TglInsert:       period.TglInsert,
		TglUpdate:       period.TglUpdate,
		UserUpdate:      period.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create Tanggal Pendaftaran
// @Description Create new registration period and set as active (deactivate others)
// @Tags Tanggal Pendaftaran
// @Accept json
// @Produce json
// @Param body body request.CreateTanggalPendaftaranRequest true "Tanggal Pendaftaran Data"
// @Success 201 {object} object{success=bool,message=string,data=response.TanggalPendaftaranResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tanggal-pendaftaran [post]
func (ctrl *TanggalPendaftaranController) Create(c *fiber.Ctx) error {
	var req request.CreateTanggalPendaftaranRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.TanggalMulai == "" || req.TanggalSelesai == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Parse dates
	tanggalMulai, err := time.Parse("2006-01-02 15:04:05", req.TanggalMulai)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tanggal_mulai format. Use: YYYY-MM-DD HH:MM:SS")
	}

	tanggalSelesai, err := time.Parse("2006-01-02 15:04:05", req.TanggalSelesai)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tanggal_selesai format. Use: YYYY-MM-DD HH:MM:SS")
	}

	// Validate date range
	if tanggalSelesai.Before(tanggalMulai) {
		return utils.BadRequestResponse(c, "Tanggal selesai must be after tanggal mulai")
	}

	// Deactivate all other periods
	database.DB.Model(&models.TanggalPendaftaran{}).
		Where("is_active = ?", 1).
		Update("is_active", 0)

	// Create new period (automatically set as active)
	now := time.Now()
	period := models.TanggalPendaftaran{
		TanggalMulai:   tanggalMulai,
		TanggalSelesai: tanggalSelesai,
		Keterangan:     req.Keterangan,
		IsActive:       1, // Always set new period as active
		Status:         req.Status,
		Hapus:          0,
		TglInsert:      &now,
		UserUpdate:     strconv.Itoa(int(utils.GetCurrentUserID(c))),
	}

	if err := database.DB.Create(&period).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create tanggal pendaftaran")
	}

	statusText := "Aktif"
	if period.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.TanggalPendaftaranResponse{
		ID:              period.ID,
		TanggalMulai:    period.TanggalMulai,
		TanggalSelesai:  period.TanggalSelesai,
		Keterangan:      period.Keterangan,
		IsActive:        period.IsActive,
		IsActiveText:    "Aktif",
		IsOpen:          period.IsOpen(),
		Status:          period.Status,
		StatusText:      statusText,
		TglInsert:       period.TglInsert,
		TglUpdate:       period.TglUpdate,
		UserUpdate:      period.UserUpdate,
	}

	return utils.CreatedResponse(c, "Tanggal pendaftaran created and set as active", result)
}

// Update godoc
// @Summary Update Tanggal Pendaftaran
// @Description Update existing registration period
// @Tags Tanggal Pendaftaran
// @Accept json
// @Produce json
// @Param id path int true "Tanggal Pendaftaran ID"
// @Param body body request.UpdateTanggalPendaftaranRequest true "Tanggal Pendaftaran Data"
// @Success 200 {object} object{success=bool,message=string,data=response.TanggalPendaftaranResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tanggal-pendaftaran/{id} [put]
func (ctrl *TanggalPendaftaranController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateTanggalPendaftaranRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.TanggalMulai == "" || req.TanggalSelesai == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Find existing
	var period models.TanggalPendaftaran
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&period).Error; err != nil {
		return utils.NotFoundResponse(c, "Tanggal pendaftaran not found")
	}

	// Parse dates
	tanggalMulai, err := time.Parse("2006-01-02 15:04:05", req.TanggalMulai)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tanggal_mulai format. Use: YYYY-MM-DD HH:MM:SS")
	}

	tanggalSelesai, err := time.Parse("2006-01-02 15:04:05", req.TanggalSelesai)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid tanggal_selesai format. Use: YYYY-MM-DD HH:MM:SS")
	}

	// Validate date range
	if tanggalSelesai.Before(tanggalMulai) {
		return utils.BadRequestResponse(c, "Tanggal selesai must be after tanggal mulai")
	}

	// Update
	period.TanggalMulai = tanggalMulai
	period.TanggalSelesai = tanggalSelesai
	period.Keterangan = req.Keterangan
	period.Status = req.Status
	period.UserUpdate = strconv.Itoa(int(utils.GetCurrentUserID(c)))

	if err := database.DB.Save(&period).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update tanggal pendaftaran")
	}

	statusText := "Aktif"
	if period.Status == 2 {
		statusText = "Tidak Aktif"
	}

	isActiveText := "Tidak Aktif"
	if period.IsActive == 1 {
		isActiveText = "Aktif"
	}

	result := response.TanggalPendaftaranResponse{
		ID:              period.ID,
		TanggalMulai:    period.TanggalMulai,
		TanggalSelesai:  period.TanggalSelesai,
		Keterangan:      period.Keterangan,
		IsActive:        period.IsActive,
		IsActiveText:    isActiveText,
		IsOpen:          period.IsOpen(),
		Status:          period.Status,
		StatusText:      statusText,
		TglInsert:       period.TglInsert,
		TglUpdate:       period.TglUpdate,
		UserUpdate:      period.UserUpdate,
	}

	return utils.SuccessResponse(c, "Tanggal pendaftaran updated successfully", result)
}

// Delete godoc
// @Summary Delete Tanggal Pendaftaran
// @Description Soft delete registration period
// @Tags Tanggal Pendaftaran
// @Accept json
// @Produce json
// @Param id path int true "Tanggal Pendaftaran ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /tanggal-pendaftaran/{id} [delete]
func (ctrl *TanggalPendaftaranController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var period models.TanggalPendaftaran
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&period).Error; err != nil {
		return utils.NotFoundResponse(c, "Tanggal pendaftaran not found")
	}

	// Soft delete
	period.Hapus = 1
	period.IsActive = 0 // Also deactivate
	period.UserUpdate = strconv.Itoa(int(utils.GetCurrentUserID(c)))

	if err := database.DB.Save(&period).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete tanggal pendaftaran")
	}

	return utils.SuccessResponse(c, "Tanggal pendaftaran deleted successfully", nil)
}