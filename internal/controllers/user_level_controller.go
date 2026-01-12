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

type UserLevelController struct{}

func NewUserLevelController() *UserLevelController {
	return &UserLevelController{}
}

// GetList godoc
// @Summary List User Levels
// @Description Get list of user levels with pagination
// @Tags User Level
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by nama_level"
// @Success 200 {object} object{success=bool,message=string,data=response.UserLevelListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /user-levels [get]
func (ctrl *UserLevelController) GetList(c *fiber.Ctx) error {
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
	query := database.DB.Model(&models.UserLevel{}).Where("hapus = ?", 0)

	// Search
	if search != "" {
		query = query.Where("nama_level LIKE ?", "%"+search+"%")
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data
	var userLevels []models.UserLevel
	if err := query.Order("id ASC").Offset(offset).Limit(perPage).Find(&userLevels).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Transform to response
	var data []response.UserLevelResponse
	for _, ul := range userLevels {
		statusText := "Aktif"
		if ul.Status == 2 {
			statusText = "Tidak Aktif"
		}

		data = append(data, response.UserLevelResponse{
			ID:         ul.ID,
			NamaLevel:  ul.NamaLevel,
			Status:     ul.Status,
			StatusText: statusText,
			TglInsert:  ul.TglInsert,
			TglUpdate:  ul.TglUpdate,
			UserUpdate: ul.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.UserLevelListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetByID godoc
// @Summary Get User Level by ID
// @Description Get user level detail by ID
// @Tags User Level
// @Accept json
// @Produce json
// @Param id path int true "User Level ID"
// @Success 200 {object} object{success=bool,message=string,data=response.UserLevelResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /user-levels/{id} [get]
func (ctrl *UserLevelController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var userLevel models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&userLevel).Error; err != nil {
		return utils.NotFoundResponse(c, "User level not found")
	}

	statusText := "Aktif"
	if userLevel.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.UserLevelResponse{
		ID:         userLevel.ID,
		NamaLevel:  userLevel.NamaLevel,
		Status:     userLevel.Status,
		StatusText: statusText,
		TglInsert:  userLevel.TglInsert,
		TglUpdate:  userLevel.TglUpdate,
		UserUpdate: userLevel.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create User Level
// @Description Create new user level
// @Tags User Level
// @Accept json
// @Produce json
// @Param body body request.CreateUserLevelRequest true "User Level Data"
// @Success 201 {object} object{success=bool,message=string,data=response.UserLevelResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /user-levels [post]
func (ctrl *UserLevelController) Create(c *fiber.Ctx) error {
	var req request.CreateUserLevelRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaLevel == "" {
		return utils.BadRequestResponse(c, "Nama level is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Check duplicate
	var count int64
	database.DB.Model(&models.UserLevel{}).Where("nama_level = ? AND hapus = ?", req.NamaLevel, 0).Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "User level with this name already exists")
	}

	// Create
	now := time.Now()
	userLevel := models.UserLevel{
		NamaLevel:    req.NamaLevel,
		Status:       req.Status,
		Hapus:        0,
		TglInsert:    &now,
		UserUpdate: strconv.Itoa(int(utils.GetCurrentUserID(c))),
	}

	if err := database.DB.Create(&userLevel).Error; err != nil {
		// Print error untuk debugging
		println("ERROR CREATE:", err.Error())
		return utils.InternalServerErrorResponse(c, "Failed to create user level: "+err.Error())
	}

	statusText := "Aktif"
	if userLevel.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.UserLevelResponse{
		ID:         userLevel.ID,
		NamaLevel:  userLevel.NamaLevel,
		Status:     userLevel.Status,
		StatusText: statusText,
		TglInsert:  userLevel.TglInsert,
		TglUpdate:  userLevel.TglUpdate,
		UserUpdate: userLevel.UserUpdate,
	}

	return utils.CreatedResponse(c, "User level created successfully", result)
}

// Update godoc
// @Summary Update User Level
// @Description Update existing user level
// @Tags User Level
// @Accept json
// @Produce json
// @Param id path int true "User Level ID"
// @Param body body request.UpdateUserLevelRequest true "User Level Data"
// @Success 200 {object} object{success=bool,message=string,data=response.UserLevelResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /user-levels/{id} [put]
func (ctrl *UserLevelController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateUserLevelRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaLevel == "" {
		return utils.BadRequestResponse(c, "Nama level is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Find existing
	var userLevel models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&userLevel).Error; err != nil {
		return utils.NotFoundResponse(c, "User level not found")
	}

	// Check duplicate (exclude current)
	var count int64
	database.DB.Model(&models.UserLevel{}).
		Where("nama_level = ? AND id != ? AND hapus = ?", req.NamaLevel, id, 0).
		Count(&count)
	if count > 0 {
		return utils.BadRequestResponse(c, "User level with this name already exists")
	}

	// Update
	userLevel.NamaLevel = req.NamaLevel
	userLevel.Status = req.Status
	userLevel.UserUpdate = "1" // User ID as string (TODO: Get from JWT token)

	if err := database.DB.Save(&userLevel).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update user level")
	}

	statusText := "Aktif"
	if userLevel.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.UserLevelResponse{
		ID:         userLevel.ID,
		NamaLevel:  userLevel.NamaLevel,
		Status:     userLevel.Status,
		StatusText: statusText,
		TglInsert:  userLevel.TglInsert,
		TglUpdate:  userLevel.TglUpdate,
		UserUpdate: userLevel.UserUpdate,
	}

	return utils.SuccessResponse(c, "User level updated successfully", result)
}

// Delete godoc
// @Summary Delete User Level
// @Description Soft delete user level
// @Tags User Level
// @Accept json
// @Produce json
// @Param id path int true "User Level ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /user-levels/{id} [delete]
func (ctrl *UserLevelController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var userLevel models.UserLevel
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&userLevel).Error; err != nil {
		return utils.NotFoundResponse(c, "User level not found")
	}

	// Soft delete
	userLevel.Hapus = 1
	userLevel.UserUpdate = "1" // User ID as string (TODO: Get from JWT token)

	if err := database.DB.Save(&userLevel).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete user level")
	}

	return utils.SuccessResponse(c, "User level deleted successfully", nil)
}