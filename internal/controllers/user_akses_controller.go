package controllers

import (
	"strconv"

	"rires-be/internal/dto/request"
	"rires-be/internal/dto/response"
	"rires-be/pkg/services"
	"rires-be/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UserAksesController handles user access management endpoints
type UserAksesController struct {
	service   *services.UserAksesService
	validator *validator.Validate
}

// NewUserAksesController creates a new controller instance
func NewUserAksesController() *UserAksesController {
	return &UserAksesController{
		service:   services.NewUserAksesService(),
		validator: validator.New(),
	}
}

// GetAllAccesses godoc
// @Summary Get All User Accesses
// @Description Admin gets all user accesses with optional filters
// @Tags Admin - User Access
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id_user_level query int false "Filter by user level"
// @Param id_menu query int false "Filter by menu"
// @Success 200 {object} response.APIResponse{data=[]response.UserAksesResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/user-akses [get]
func (ctrl *UserAksesController) GetAllAccesses(c *fiber.Ctx) error {
	// Parse filters
	idUserLevel, _ := strconv.Atoi(c.Query("id_user_level", "0"))
	idMenu, _ := strconv.Atoi(c.Query("id_menu", "0"))

	// Call service
	result, err := ctrl.service.GetAllAccesses(idUserLevel, idMenu)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"Failed to get accesses",
			err.Error(),
		))
	}

	return c.JSON(response.SuccessResponse(
		"Accesses retrieved successfully",
		result,
	))
}

// GetAccessesByUserLevel godoc
// @Summary Get Accesses by User Level
// @Description Admin gets all accesses for a specific user level
// @Tags Admin - User Access
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id_user_level path int true "User Level ID"
// @Success 200 {object} response.APIResponse{data=response.UserAksesGroupedResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/user-akses/by-level/{id_user_level} [get]
func (ctrl *UserAksesController) GetAccessesByUserLevel(c *fiber.Ctx) error {
	// Parse ID
	idUserLevel, err := strconv.Atoi(c.Params("id_user_level"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid user level ID",
			err.Error(),
		))
	}

	// Call service
	result, err := ctrl.service.GetAccessesByUserLevel(idUserLevel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(
			"Failed to get accesses",
			err.Error(),
		))
	}

	return c.JSON(response.SuccessResponse(
		"Accesses retrieved successfully",
		result,
	))
}

// GetAccessDetail godoc
// @Summary Get Access Detail
// @Description Admin gets single access detail
// @Tags Admin - User Access
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Access ID"
// @Success 200 {object} response.APIResponse{data=response.UserAksesResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/user-akses/{id} [get]
func (ctrl *UserAksesController) GetAccessDetail(c *fiber.Ctx) error {
	// Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid access ID",
			err.Error(),
		))
	}

	// Call service
	result, err := ctrl.service.GetAccessDetail(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(
			"Access not found",
			err.Error(),
		))
	}

	return c.JSON(response.SuccessResponse(
		"Access detail",
		result,
	))
}

// CreateAccess godoc
// @Summary Create User Access
// @Description Admin creates new user access
// @Tags Admin - User Access
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body request.CreateUserAksesRequest true "Access data"
// @Success 201 {object} response.APIResponse{data=response.UserAksesResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/user-akses [post]
func (ctrl *UserAksesController) CreateAccess(c *fiber.Ctx) error {
	// 1. Parse request body
	var req request.CreateUserAksesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid request body",
			err.Error(),
		))
	}

	// 2. Validate request
	if err := ctrl.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Validation failed",
			err.Error(),
		))
	}

	// 3. Get user ID
	userID := int(utils.GetCurrentUserID(c))

	// 4. Call service
	result, err := ctrl.service.CreateAccess(&req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to create access",
			err.Error(),
		))
	}

	// 5. Return success
	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Akses berhasil dibuat",
		result,
	))
}

// UpdateAccess godoc
// @Summary Update User Access
// @Description Admin updates existing access
// @Tags Admin - User Access
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Access ID"
// @Param body body request.UpdateUserAksesRequest true "Update data"
// @Success 200 {object} response.APIResponse{data=response.UserAksesResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/user-akses/{id} [put]
func (ctrl *UserAksesController) UpdateAccess(c *fiber.Ctx) error {
	// 1. Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid access ID",
			err.Error(),
		))
	}

	// 2. Parse request body
	var req request.UpdateUserAksesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid request body",
			err.Error(),
		))
	}

	// 3. Validate request
	if err := ctrl.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Validation failed",
			err.Error(),
		))
	}

	// 4. Get user ID
	userID := int(utils.GetCurrentUserID(c))

	// 5. Call service
	result, err := ctrl.service.UpdateAccess(id, &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to update access",
			err.Error(),
		))
	}

	// 6. Return success
	return c.JSON(response.SuccessResponse(
		"Akses berhasil diupdate",
		result,
	))
}

// DeleteAccess godoc
// @Summary Delete User Access
// @Description Admin deletes access
// @Tags Admin - User Access
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Access ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/user-akses/{id} [delete]
func (ctrl *UserAksesController) DeleteAccess(c *fiber.Ctx) error {
	// 1. Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid access ID",
			err.Error(),
		))
	}

	// 2. Get user ID
	userID := int(utils.GetCurrentUserID(c))

	// 3. Call service
	if err := ctrl.service.DeleteAccess(id, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to delete access",
			err.Error(),
		))
	}

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Akses berhasil dihapus",
		nil,
	))
}

// BulkCreateAccess godoc
// @Summary Bulk Create User Accesses
// @Description Admin creates multiple accesses at once for a user level
// @Tags Admin - User Access
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body request.BulkCreateUserAksesRequest true "Bulk create data"
// @Success 201 {object} response.APIResponse{data=[]response.UserAksesResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/user-akses/bulk [post]
func (ctrl *UserAksesController) BulkCreateAccess(c *fiber.Ctx) error {
	// 1. Parse request body
	var req request.BulkCreateUserAksesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid request body",
			err.Error(),
		))
	}

	// 2. Validate request
	if err := ctrl.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Validation failed",
			err.Error(),
		))
	}

	// 3. Get user ID
	userID := int(utils.GetCurrentUserID(c))

	// 4. Call service
	result, err := ctrl.service.BulkCreateAccess(&req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to create accesses",
			err.Error(),
		))
	}

	// 5. Return success
	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Akses berhasil dibuat secara bulk",
		result,
	))
}