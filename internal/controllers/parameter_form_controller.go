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

type ParameterFormController struct{}

func NewParameterFormController() *ParameterFormController {
	return &ParameterFormController{}
}

// GetList godoc
// @Summary List Parameter Form
// @Description Get list of parameter form with pagination
// @Tags Parameter Form
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by nama_parameter or label"
// @Param kategori_id query int false "Filter by kategori_id"
// @Success 200 {object} object{success=bool,message=string,data=response.ParameterFormListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /parameter-form [get]
func (ctrl *ParameterFormController) GetList(c *fiber.Ctx) error {
	// Pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	kategoriID := c.Query("kategori_id", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Query builder
	query := database.DB.Model(&models.ParameterForm{}).Where("hapus = ?", 0)

	// Filter by kategori
	if kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}

	// Search
	if search != "" {
		query = query.Where("nama_parameter LIKE ? OR label LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data with kategori
	var params []models.ParameterForm
	if err := query.Preload("Kategori").Order("`kategori_id` ASC, `urutan` ASC").
		Offset(offset).Limit(perPage).Find(&params).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Transform to response
	var data []response.ParameterFormResponse
	for _, param := range params {
		statusText := "Aktif"
		if param.Status == 2 {
			statusText = "Tidak Aktif"
		}

		namaKategori := ""
		if param.Kategori != nil {
			namaKategori = param.Kategori.NamaKategori
		}

		data = append(data, response.ParameterFormResponse{
			ID:            param.ID,
			KategoriID:    param.KategoriID,
			NamaKategori:  namaKategori,
			NamaParameter: param.NamaParameter,
			Label:         param.Label,
			TipeInput:     param.TipeInput,
			Validasi:      param.Validasi,
			Placeholder:   param.Placeholder,
			HelpText:      param.HelpText,
			Opsi:          param.Opsi,
			Urutan:        param.Urutan,
			Status:        param.Status,
			StatusText:    statusText,
			TglInsert:     param.TglInsert,
			TglUpdate:     param.TglUpdate,
			UserUpdate:    param.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.ParameterFormListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetByKategori godoc
// @Summary Get Parameter Form by Kategori
// @Description Get all parameters for specific kategori (for dynamic form rendering)
// @Tags Parameter Form
// @Accept json
// @Produce json
// @Param kategori_id path int true "Kategori PKM ID"
// @Success 200 {object} object{success=bool,message=string,data=response.ParameterFormByKategoriResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /parameter-form/kategori/{kategori_id} [get]
func (ctrl *ParameterFormController) GetByKategori(c *fiber.Ctx) error {
	kategoriID, err := strconv.Atoi(c.Params("kategori_id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid kategori_id")
	}

	// Check if kategori exists
	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", kategoriID, 0).First(&kategori).Error; err != nil {
		return utils.NotFoundResponse(c, "Kategori PKM not found")
	}

	// Get parameters for this kategori
	var params []models.ParameterForm
	if err := database.DB.Where("kategori_id = ? AND hapus = ? AND status = ?", kategoriID, 0, 1).
		Order("`urutan` ASC").Find(&params).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch parameters")
	}

	// Transform to response
	var paramResponses []response.ParameterFormResponse
	for _, param := range params {
		paramResponses = append(paramResponses, response.ParameterFormResponse{
			ID:            param.ID,
			KategoriID:    param.KategoriID,
			NamaParameter: param.NamaParameter,
			Label:         param.Label,
			TipeInput:     param.TipeInput,
			Validasi:      param.Validasi,
			Placeholder:   param.Placeholder,
			HelpText:      param.HelpText,
			Opsi:          param.Opsi,
			Urutan:        param.Urutan,
			Status:        param.Status,
		})
	}

	result := response.ParameterFormByKategoriResponse{
		KategoriID:   kategori.ID,
		NamaKategori: kategori.NamaKategori,
		Parameters:   paramResponses,
	}

	return utils.SuccessResponse(c, "Parameters retrieved successfully", result)
}

// GetByID godoc
// @Summary Get Parameter Form by ID
// @Description Get parameter form detail by ID
// @Tags Parameter Form
// @Accept json
// @Produce json
// @Param id path int true "Parameter Form ID"
// @Success 200 {object} object{success=bool,message=string,data=response.ParameterFormResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /parameter-form/{id} [get]
func (ctrl *ParameterFormController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var param models.ParameterForm
	if err := database.DB.Preload("Kategori").Where("id = ? AND hapus = ?", id, 0).First(&param).Error; err != nil {
		return utils.NotFoundResponse(c, "Parameter form not found")
	}

	statusText := "Aktif"
	if param.Status == 2 {
		statusText = "Tidak Aktif"
	}

	namaKategori := ""
	if param.Kategori != nil {
		namaKategori = param.Kategori.NamaKategori
	}

	result := response.ParameterFormResponse{
		ID:            param.ID,
		KategoriID:    param.KategoriID,
		NamaKategori:  namaKategori,
		NamaParameter: param.NamaParameter,
		Label:         param.Label,
		TipeInput:     param.TipeInput,
		Validasi:      param.Validasi,
		Placeholder:   param.Placeholder,
		HelpText:      param.HelpText,
		Opsi:          param.Opsi,
		Urutan:        param.Urutan,
		Status:        param.Status,
		StatusText:    statusText,
		TglInsert:     param.TglInsert,
		TglUpdate:     param.TglUpdate,
		UserUpdate:    param.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create Parameter Form
// @Description Create new parameter form
// @Tags Parameter Form
// @Accept json
// @Produce json
// @Param body body request.CreateParameterFormRequest true "Parameter Form Data"
// @Success 201 {object} object{success=bool,message=string,data=response.ParameterFormResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /parameter-form [post]
func (ctrl *ParameterFormController) Create(c *fiber.Ctx) error {
	var req request.CreateParameterFormRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaParameter == "" || req.Label == "" || req.TipeInput == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Check if kategori exists
	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", req.KategoriID, 0).First(&kategori).Error; err != nil {
		return utils.BadRequestResponse(c, "Kategori PKM not found")
	}

	// Create
	now := time.Now()
	param := models.ParameterForm{
		KategoriID:    req.KategoriID,
		NamaParameter: req.NamaParameter,
		Label:         req.Label,
		TipeInput:     req.TipeInput,
		Validasi:      req.Validasi,
		Placeholder:   req.Placeholder,
		HelpText:      req.HelpText,
		Opsi:          req.Opsi,
		Urutan:        req.Urutan,
		Status:        req.Status,
		Hapus:         0,
		TglInsert:     &now,
		UserUpdate:    "1", // TODO: Get from JWT token
	}

	if err := database.DB.Create(&param).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create parameter form")
	}

	statusText := "Aktif"
	if param.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.ParameterFormResponse{
		ID:            param.ID,
		KategoriID:    param.KategoriID,
		NamaKategori:  kategori.NamaKategori,
		NamaParameter: param.NamaParameter,
		Label:         param.Label,
		TipeInput:     param.TipeInput,
		Validasi:      param.Validasi,
		Placeholder:   param.Placeholder,
		HelpText:      param.HelpText,
		Urutan:        param.Urutan,
		Status:        param.Status,
		StatusText:    statusText,
		TglInsert:     param.TglInsert,
		TglUpdate:     param.TglUpdate,
		UserUpdate:    param.UserUpdate,
	}

	return utils.CreatedResponse(c, "Parameter form created successfully", result)
}

// Update godoc
// @Summary Update Parameter Form
// @Description Update existing parameter form
// @Tags Parameter Form
// @Accept json
// @Produce json
// @Param id path int true "Parameter Form ID"
// @Param body body request.UpdateParameterFormRequest true "Parameter Form Data"
// @Success 200 {object} object{success=bool,message=string,data=response.ParameterFormResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /parameter-form/{id} [put]
func (ctrl *ParameterFormController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateParameterFormRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaParameter == "" || req.Label == "" || req.TipeInput == "" {
		return utils.BadRequestResponse(c, "Required fields are missing")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Find existing
	var param models.ParameterForm
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&param).Error; err != nil {
		return utils.NotFoundResponse(c, "Parameter form not found")
	}

	// Check if kategori exists
	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", req.KategoriID, 0).First(&kategori).Error; err != nil {
		return utils.BadRequestResponse(c, "Kategori PKM not found")
	}

	// Update
	param.KategoriID = req.KategoriID
	param.NamaParameter = req.NamaParameter
	param.Label = req.Label
	param.TipeInput = req.TipeInput
	param.Validasi = req.Validasi
	param.Placeholder = req.Placeholder
	param.HelpText = req.HelpText
	param.Opsi = req.Opsi
	param.Urutan = req.Urutan
	param.Status = req.Status
	param.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&param).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update parameter form")
	}

	statusText := "Aktif"
	if param.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.ParameterFormResponse{
		ID:            param.ID,
		KategoriID:    param.KategoriID,
		NamaKategori:  kategori.NamaKategori,
		NamaParameter: param.NamaParameter,
		Label:         param.Label,
		TipeInput:     param.TipeInput,
		Validasi:      param.Validasi,
		Placeholder:   param.Placeholder,
		HelpText:      param.HelpText,
		Urutan:        param.Urutan,
		Status:        param.Status,
		StatusText:    statusText,
		TglInsert:     param.TglInsert,
		TglUpdate:     param.TglUpdate,
		UserUpdate:    param.UserUpdate,
	}

	return utils.SuccessResponse(c, "Parameter form updated successfully", result)
}

// Delete godoc
// @Summary Delete Parameter Form
// @Description Soft delete parameter form
// @Tags Parameter Form
// @Accept json
// @Produce json
// @Param id path int true "Parameter Form ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /parameter-form/{id} [delete]
func (ctrl *ParameterFormController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var param models.ParameterForm
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&param).Error; err != nil {
		return utils.NotFoundResponse(c, "Parameter form not found")
	}

	// Soft delete
	param.Hapus = 1
	param.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&param).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete parameter form")
	}

	return utils.SuccessResponse(c, "Parameter form deleted successfully", nil)
}