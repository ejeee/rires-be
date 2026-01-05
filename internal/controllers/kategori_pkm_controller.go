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

type KategoriPKMController struct{}

func NewKategoriPKMController() *KategoriPKMController {
	return &KategoriPKMController{}
}

// GetList godoc
// @Summary List Kategori PKM
// @Description Get list of kategori PKM with pagination
// @Tags Kategori PKM
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by nama_kategori"
// @Success 200 {object} object{success=bool,message=string,data=response.KategoriPKMListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /kategori-pkm [get]
func (ctrl *KategoriPKMController) GetList(c *fiber.Ctx) error {
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
	query := database.DB.Model(&models.KategoriPKM{}).Where("hapus = ?", 0)

	// Search
	if search != "" {
		query = query.Where("nama_kategori LIKE ?", "%"+search+"%")
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data
	var kategoris []models.KategoriPKM
	if err := query.Order("id ASC").Offset(offset).Limit(perPage).Find(&kategoris).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Transform to response
	var data []response.KategoriPKMResponse
	for _, kat := range kategoris {
		statusText := "Aktif"
		if kat.Status == 2 {
			statusText = "Tidak Aktif"
		}

		data = append(data, response.KategoriPKMResponse{
			ID:           kat.ID,
			NamaKategori: kat.NamaKategori,
			Status:       kat.Status,
			StatusText:   statusText,
			TglInsert:    kat.TglInsert,
			TglUpdate:    kat.TglUpdate,
			UserUpdate:   kat.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.KategoriPKMListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetByID godoc
// @Summary Get Kategori PKM by ID
// @Description Get kategori PKM detail by ID
// @Tags Kategori PKM
// @Accept json
// @Produce json
// @Param id path int true "Kategori PKM ID"
// @Success 200 {object} object{success=bool,message=string,data=response.KategoriPKMResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /kategori-pkm/{id} [get]
func (ctrl *KategoriPKMController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&kategori).Error; err != nil {
		return utils.NotFoundResponse(c, "Kategori PKM not found")
	}

	statusText := "Aktif"
	if kategori.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.KategoriPKMResponse{
		ID:           kategori.ID,
		NamaKategori: kategori.NamaKategori,
		Status:       kategori.Status,
		StatusText:   statusText,
		TglInsert:    kategori.TglInsert,
		TglUpdate:    kategori.TglUpdate,
		UserUpdate:   kategori.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create Kategori PKM
// @Description Create new kategori PKM
// @Tags Kategori PKM
// @Accept json
// @Produce json
// @Param body body request.CreateKategoriPKMRequest true "Kategori PKM Data"
// @Success 201 {object} object{success=bool,message=string,data=response.KategoriPKMResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /kategori-pkm [post]
func (ctrl *KategoriPKMController) Create(c *fiber.Ctx) error {
	var req request.CreateKategoriPKMRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaKategori == "" {
		return utils.BadRequestResponse(c, "Nama kategori is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Check duplicate
	var count int64
	database.DB.Model(&models.KategoriPKM{}).Where("nama_kategori = ? AND hapus = ?", req.NamaKategori, 0).Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "Kategori PKM with this name already exists")
	}

	// Create
	now := time.Now()
	kategori := models.KategoriPKM{
		NamaKategori: req.NamaKategori,
		Status:       req.Status,
		Hapus:        0,
		TglInsert:    &now,
		UserUpdate:   "1", // TODO: Get from JWT token
	}

	if err := database.DB.Create(&kategori).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create kategori PKM")
	}

	statusText := "Aktif"
	if kategori.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.KategoriPKMResponse{
		ID:           kategori.ID,
		NamaKategori: kategori.NamaKategori,
		Status:       kategori.Status,
		StatusText:   statusText,
		TglInsert:    kategori.TglInsert,
		TglUpdate:    kategori.TglUpdate,
		UserUpdate:   kategori.UserUpdate,
	}

	return utils.CreatedResponse(c, "Kategori PKM created successfully", result)
}

// Update godoc
// @Summary Update Kategori PKM
// @Description Update existing kategori PKM
// @Tags Kategori PKM
// @Accept json
// @Produce json
// @Param id path int true "Kategori PKM ID"
// @Param body body request.UpdateKategoriPKMRequest true "Kategori PKM Data"
// @Success 200 {object} object{success=bool,message=string,data=response.KategoriPKMResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /kategori-pkm/{id} [put]
func (ctrl *KategoriPKMController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateKategoriPKMRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaKategori == "" {
		return utils.BadRequestResponse(c, "Nama kategori is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Find existing
	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&kategori).Error; err != nil {
		return utils.NotFoundResponse(c, "Kategori PKM not found")
	}

	// Check duplicate (exclude current)
	var count int64
	database.DB.Model(&models.KategoriPKM{}).
		Where("nama_kategori = ? AND id != ? AND hapus = ?", req.NamaKategori, id, 0).
		Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "Kategori PKM with this name already exists")
	}

	// Update
	kategori.NamaKategori = req.NamaKategori
	kategori.Status = req.Status
	kategori.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&kategori).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update kategori PKM")
	}

	statusText := "Aktif"
	if kategori.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.KategoriPKMResponse{
		ID:           kategori.ID,
		NamaKategori: kategori.NamaKategori,
		Status:       kategori.Status,
		StatusText:   statusText,
		TglInsert:    kategori.TglInsert,
		TglUpdate:    kategori.TglUpdate,
		UserUpdate:   kategori.UserUpdate,
	}

	return utils.SuccessResponse(c, "Kategori PKM updated successfully", result)
}

// Delete godoc
// @Summary Delete Kategori PKM
// @Description Soft delete kategori PKM
// @Tags Kategori PKM
// @Accept json
// @Produce json
// @Param id path int true "Kategori PKM ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /kategori-pkm/{id} [delete]
func (ctrl *KategoriPKMController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&kategori).Error; err != nil {
		return utils.NotFoundResponse(c, "Kategori PKM not found")
	}

	// TODO: Check if kategori is used in any pengajuan
	// var pengajuanCount int64
	// database.DB.Model(&models.Pengajuan{}).Where("kategori_id = ?", id).Count(&pengajuanCount)
	// if pengajuanCount > 0 {
	//     return utils.BadRequestResponse(c, "Cannot delete kategori that is used in pengajuan")
	// }

	// Soft delete
	kategori.Hapus = 1
	kategori.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&kategori).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete kategori PKM")
	}

	return utils.SuccessResponse(c, "Kategori PKM deleted successfully", nil)
}