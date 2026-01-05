package controllers

import (
	"math"
	"rires-be/internal/dto/request"
	"rires-be/internal/dto/response"
	"rires-be/internal/models"
	"rires-be/pkg/database"
	"rires-be/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type StatusReviewController struct{}

func NewStatusReviewController() *StatusReviewController {
	return &StatusReviewController{}
}

// GetList godoc
// @Summary List Status Review
// @Description Get list of status review with pagination
// @Tags Status Review
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by nama_status or kode_status"
// @Success 200 {object} object{success=bool,message=string,data=response.StatusReviewListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /status-review [get]
func (ctrl *StatusReviewController) GetList(c *fiber.Ctx) error {
	// Pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Query builder
	query := database.DB.Model(&models.StatusReview{}).Where("hapus = ?", 0)

	// Search
	if search != "" {
		query = query.Where("nama_status LIKE ? OR kode_status LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data
	var statusReviews []models.StatusReview
	if err := query.Order("`urutan` ASC").Offset(offset).Limit(perPage).Find(&statusReviews).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Transform to response
	var data []response.StatusReviewResponse
	for _, sr := range statusReviews {
		statusText := "Aktif"
		if sr.Status == 2 {
			statusText = "Tidak Aktif"
		}

		data = append(data, response.StatusReviewResponse{
			ID:         sr.ID,
			NamaStatus: sr.NamaStatus,
			KodeStatus: sr.KodeStatus,
			Warna:      sr.Warna,
			Urutan:     sr.Urutan,
			Status:     sr.Status,
			StatusText: statusText,
			TglInsert:  sr.TglInsert,
			TglUpdate:  sr.TglUpdate,
			UserUpdate: sr.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.StatusReviewListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetByID godoc
// @Summary Get Status Review by ID
// @Description Get status review detail by ID
// @Tags Status Review
// @Accept json
// @Produce json
// @Param id path int true "Status Review ID"
// @Success 200 {object} object{success=bool,message=string,data=response.StatusReviewResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /status-review/{id} [get]
func (ctrl *StatusReviewController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var statusReview models.StatusReview
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&statusReview).Error; err != nil {
		return utils.NotFoundResponse(c, "Status review not found")
	}

	statusText := "Aktif"
	if statusReview.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.StatusReviewResponse{
		ID:         statusReview.ID,
		NamaStatus: statusReview.NamaStatus,
		KodeStatus: statusReview.KodeStatus,
		Warna:      statusReview.Warna,
		Urutan:     statusReview.Urutan,
		Status:     statusReview.Status,
		StatusText: statusText,
		TglInsert:  statusReview.TglInsert,
		TglUpdate:  statusReview.TglUpdate,
		UserUpdate: statusReview.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create Status Review
// @Description Create new status review
// @Tags Status Review
// @Accept json
// @Produce json
// @Param body body request.CreateStatusReviewRequest true "Status Review Data"
// @Success 201 {object} object{success=bool,message=string,data=response.StatusReviewResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /status-review [post]
func (ctrl *StatusReviewController) Create(c *fiber.Ctx) error {
	var req request.CreateStatusReviewRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaStatus == "" {
		return utils.BadRequestResponse(c, "Nama status is required")
	}
	if req.KodeStatus == "" {
		return utils.BadRequestResponse(c, "Kode status is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Uppercase kode_status
	req.KodeStatus = strings.ToUpper(req.KodeStatus)

	// Check duplicate kode
	var count int64
	database.DB.Model(&models.StatusReview{}).Where("kode_status = ? AND hapus = ?", req.KodeStatus, 0).Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "Status review with this kode already exists")
	}

	// Create
	now := time.Now()
	statusReview := models.StatusReview{
		NamaStatus: req.NamaStatus,
		KodeStatus: req.KodeStatus,
		Warna:      req.Warna,
		Urutan:     req.Urutan,
		Status:     req.Status,
		Hapus:      0,
		TglInsert:  &now,
		UserUpdate: "1", // TODO: Get from JWT token
	}

	if err := database.DB.Create(&statusReview).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create status review")
	}

	statusText := "Aktif"
	if statusReview.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.StatusReviewResponse{
		ID:         statusReview.ID,
		NamaStatus: statusReview.NamaStatus,
		KodeStatus: statusReview.KodeStatus,
		Warna:      statusReview.Warna,
		Urutan:     statusReview.Urutan,
		Status:     statusReview.Status,
		StatusText: statusText,
		TglInsert:  statusReview.TglInsert,
		TglUpdate:  statusReview.TglUpdate,
		UserUpdate: statusReview.UserUpdate,
	}

	return utils.CreatedResponse(c, "Status review created successfully", result)
}

// Update godoc
// @Summary Update Status Review
// @Description Update existing status review
// @Tags Status Review
// @Accept json
// @Produce json
// @Param id path int true "Status Review ID"
// @Param body body request.UpdateStatusReviewRequest true "Status Review Data"
// @Success 200 {object} object{success=bool,message=string,data=response.StatusReviewResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /status-review/{id} [put]
func (ctrl *StatusReviewController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateStatusReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaStatus == "" {
		return utils.BadRequestResponse(c, "Nama status is required")
	}
	if req.KodeStatus == "" {
		return utils.BadRequestResponse(c, "Kode status is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Uppercase kode_status
	req.KodeStatus = strings.ToUpper(req.KodeStatus)

	// Find existing
	var statusReview models.StatusReview
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&statusReview).Error; err != nil {
		return utils.NotFoundResponse(c, "Status review not found")
	}

	// Check duplicate kode (exclude current)
	var count int64
	database.DB.Model(&models.StatusReview{}).
		Where("kode_status = ? AND id != ? AND hapus = ?", req.KodeStatus, id, 0).
		Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "Status review with this kode already exists")
	}

	// Update
	statusReview.NamaStatus = req.NamaStatus
	statusReview.KodeStatus = req.KodeStatus
	statusReview.Warna = req.Warna
	statusReview.Urutan = req.Urutan
	statusReview.Status = req.Status
	statusReview.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&statusReview).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update status review")
	}

	statusText := "Aktif"
	if statusReview.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.StatusReviewResponse{
		ID:         statusReview.ID,
		NamaStatus: statusReview.NamaStatus,
		KodeStatus: statusReview.KodeStatus,
		Warna:      statusReview.Warna,
		Urutan:     statusReview.Urutan,
		Status:     statusReview.Status,
		StatusText: statusText,
		TglInsert:  statusReview.TglInsert,
		TglUpdate:  statusReview.TglUpdate,
		UserUpdate: statusReview.UserUpdate,
	}

	return utils.SuccessResponse(c, "Status review updated successfully", result)
}

// Delete godoc
// @Summary Delete Status Review
// @Description Soft delete status review
// @Tags Status Review
// @Accept json
// @Produce json
// @Param id path int true "Status Review ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /status-review/{id} [delete]
func (ctrl *StatusReviewController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var statusReview models.StatusReview
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&statusReview).Error; err != nil {
		return utils.NotFoundResponse(c, "Status review not found")
	}

	// TODO: Check if status is used in any review
	// var reviewCount int64
	// database.DB.Model(&models.Review{}).Where("status_id = ?", id).Count(&reviewCount)
	// if reviewCount > 0 {
	//     return utils.BadRequestResponse(c, "Cannot delete status that is used in reviews")
	// }

	// Soft delete
	statusReview.Hapus = 1
	statusReview.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&statusReview).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete status review")
	}

	return utils.SuccessResponse(c, "Status review deleted successfully", nil)
}