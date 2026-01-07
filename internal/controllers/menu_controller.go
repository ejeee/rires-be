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

type MenuController struct{}

func NewMenuController() *MenuController {
	return &MenuController{}
}

// GetList godoc
// @Summary List Menus
// @Description Get list of menus with pagination (flat structure)
// @Tags Menu
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search by nama_menu"
// @Param id_parent query int false "Filter by id_parent"
// @Success 200 {object} object{success=bool,message=string,data=response.MenuListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /menus [get]
func (ctrl *MenuController) GetList(c *fiber.Ctx) error {
	// Pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	IDParent := c.Query("id_parent", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Query builder
	query := database.DB.Model(&models.Menu{}).Where("hapus = ?", 0)

	// Search
	if search != "" {
		query = query.Where("nama_menu LIKE ?", "%"+search+"%")
	}

	// Filter by parent
	if IDParent != "" {
		query = query.Where("id_parent = ?", IDParent)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to count data")
	}

	// Get data
	var menus []models.Menu
	if err := query.Order("id_parent ASC, urutan ASC").Offset(offset).Limit(perPage).Find(&menus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch data",
			"error":   err.Error(), // Show actual error
		})
	}

	// Transform to response
	var data []response.MenuResponse
	for _, menu := range menus {
		statusText := "Aktif"
		if menu.Status == 2 {
			statusText = "Tidak Aktif"
		}

		data = append(data, response.MenuResponse{
			ID:         menu.ID,
			IDParent:   menu.IDParent,
			NamaMenu:   menu.NamaMenu,
			URLMenu:    menu.URLMenu,
			Lucide:     menu.Lucide,
			Urutan:     menu.Urutan,
			Status:     menu.Status,
			StatusText: statusText,
			TglInsert:  menu.TglInsert,
			TglUpdate:  menu.TglUpdate,
			UserUpdate: menu.UserUpdate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	result := response.MenuListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// GetTree godoc
// @Summary Get Menu Tree
// @Description Get menus in tree/hierarchical structure (for sidebar/navigation)
// @Tags Menu
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string,data=[]response.MenuTreeResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /menus/tree [get]
func (ctrl *MenuController) GetTree(c *fiber.Ctx) error {
	// Get all active menus
	var menus []models.Menu
	if err := database.DB.Where("hapus = ? AND status = ?", 0, 1).
		Order("id_parent ASC, urutan ASC").
		Find(&menus).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch data")
	}

	// Build tree structure
	tree := buildMenuTree(menus, 0)

	return utils.SuccessResponse(c, "Menu tree retrieved successfully", tree)
}

// Helper function to build tree structure
func buildMenuTree(menus []models.Menu, IDParent int) []response.MenuTreeResponse {
	var tree []response.MenuTreeResponse

	for _, menu := range menus {
		if menu.IDParent == IDParent {
			statusText := "Aktif"
			if menu.Status == 2 {
				statusText = "Tidak Aktif"
			}

			node := response.MenuTreeResponse{
				ID:         menu.ID,
				IDParent:   menu.IDParent,
				NamaMenu:   menu.NamaMenu,
				URLMenu:    menu.URLMenu,
				Lucide:     menu.Lucide,
				Urutan:     menu.Urutan,
				Status:     menu.Status,
				StatusText: statusText,
				Children:   buildMenuTree(menus, menu.ID),
			}

			tree = append(tree, node)
		}
	}

	return tree
}

// GetByID godoc
// @Summary Get Menu by ID
// @Description Get menu detail by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} object{success=bool,message=string,data=response.MenuResponse}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /menus/{id} [get]
func (ctrl *MenuController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var menu models.Menu
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&menu).Error; err != nil {
		return utils.NotFoundResponse(c, "Menu not found")
	}

	statusText := "Aktif"
	if menu.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.MenuResponse{
		ID:         menu.ID,
		IDParent:   menu.IDParent,
		NamaMenu:   menu.NamaMenu,
		URLMenu:    menu.URLMenu,
		Lucide:     menu.Lucide,
		Urutan:     menu.Urutan,
		Status:     menu.Status,
		StatusText: statusText,
		TglInsert:  menu.TglInsert,
		TglUpdate:  menu.TglUpdate,
		UserUpdate: menu.UserUpdate,
	}

	return utils.SuccessResponse(c, "Data retrieved successfully", result)
}

// Create godoc
// @Summary Create Menu
// @Description Create new menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param body body request.CreateMenuRequest true "Menu Data"
// @Success 201 {object} object{success=bool,message=string,data=response.MenuResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /menus [post]
func (ctrl *MenuController) Create(c *fiber.Ctx) error {
	var req request.CreateMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaMenu == "" {
		return utils.BadRequestResponse(c, "Nama menu is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Check if parent exists (if id_parent > 0)
	if req.IDParent > 0 {
		var parentMenu models.Menu
		if err := database.DB.Where("id = ? AND hapus = ?", req.IDParent, 0).First(&parentMenu).Error; err != nil {
			return utils.BadRequestResponse(c, "Parent menu not found")
		}
	}

	// Create
	now := time.Now()
	menu := models.Menu{
		IDParent:   req.IDParent,
		NamaMenu:   req.NamaMenu,
		URLMenu:    req.URLMenu,
		Lucide:     req.Lucide,
		Urutan:     req.Urutan,
		Status:     req.Status,
		Hapus:      0,
		TglInsert:  &now,
		UserUpdate: "1", // TODO: Get from JWT token
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create menu")
	}

	statusText := "Aktif"
	if menu.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.MenuResponse{
		ID:         menu.ID,
		IDParent:   menu.IDParent,
		NamaMenu:   menu.NamaMenu,
		URLMenu:    menu.URLMenu,
		Lucide:     menu.Lucide,
		Urutan:     menu.Urutan,
		Status:     menu.Status,
		StatusText: statusText,
		TglInsert:  menu.TglInsert,
		TglUpdate:  menu.TglUpdate,
		UserUpdate: menu.UserUpdate,
	}

	return utils.CreatedResponse(c, "Menu created successfully", result)
}

// Update godoc
// @Summary Update Menu
// @Description Update existing menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param body body request.UpdateMenuRequest true "Menu Data"
// @Success 200 {object} object{success=bool,message=string,data=response.MenuResponse}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /menus/{id} [put]
func (ctrl *MenuController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	var req request.UpdateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate
	if req.NamaMenu == "" {
		return utils.BadRequestResponse(c, "Nama menu is required")
	}
	if req.Status != 1 && req.Status != 2 {
		return utils.BadRequestResponse(c, "Status must be 1 (active) or 2 (inactive)")
	}

	// Find existing
	var menu models.Menu
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&menu).Error; err != nil {
		return utils.NotFoundResponse(c, "Menu not found")
	}

	// Prevent circular parent reference
	if req.IDParent == id {
		return utils.BadRequestResponse(c, "Menu cannot be its own parent")
	}

	// Check if parent exists (if id_parent > 0)
	if req.IDParent > 0 {
		var parentMenu models.Menu
		if err := database.DB.Where("id = ? AND hapus = ?", req.IDParent, 0).First(&parentMenu).Error; err != nil {
			return utils.BadRequestResponse(c, "Parent menu not found")
		}
	}

	// Update
	menu.IDParent = req.IDParent
	menu.NamaMenu = req.NamaMenu
	menu.URLMenu = req.URLMenu
	menu.Lucide = req.Lucide
	menu.Urutan = req.Urutan
	menu.Status = req.Status
	menu.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&menu).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update menu")
	}

	statusText := "Aktif"
	if menu.Status == 2 {
		statusText = "Tidak Aktif"
	}

	result := response.MenuResponse{
		ID:         menu.ID,
		IDParent:   menu.IDParent,
		NamaMenu:   menu.NamaMenu,
		URLMenu:    menu.URLMenu,
		Lucide:     menu.Lucide,
		Urutan:     menu.Urutan,
		Status:     menu.Status,
		StatusText: statusText,
		TglInsert:  menu.TglInsert,
		TglUpdate:  menu.TglUpdate,
		UserUpdate: menu.UserUpdate,
	}

	return utils.SuccessResponse(c, "Menu updated successfully", result)
}

// Delete godoc
// @Summary Delete Menu
// @Description Soft delete menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{success=bool,message=string}
// @Failure 404 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /menus/{id} [delete]
func (ctrl *MenuController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid ID")
	}

	// Find existing
	var menu models.Menu
	if err := database.DB.Where("id = ? AND hapus = ?", id, 0).First(&menu).Error; err != nil {
		return utils.NotFoundResponse(c, "Menu not found")
	}

	// Check if has children
	var childCount int64
	database.DB.Model(&models.Menu{}).Where("id_parent = ? AND hapus = ?", id, 0).Count(&childCount)
	if childCount > 0 {
		return utils.BadRequestResponse(c, "Cannot delete menu with children. Delete children first.")
	}

	// Soft delete
	menu.Hapus = 1
	menu.UserUpdate = "1" // TODO: Get from JWT token

	if err := database.DB.Save(&menu).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete menu")
	}

	return utils.SuccessResponse(c, "Menu deleted successfully", nil)
}