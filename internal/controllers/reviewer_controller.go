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

// ReviewerController handles reviewer management endpoints
type ReviewerController struct {
	service   *services.ReviewerService
	validator *validator.Validate
}

// NewReviewerController creates a new controller instance
func NewReviewerController() *ReviewerController {
	return &ReviewerController{
		service:   services.NewReviewerService(),
		validator: validator.New(),
	}
}

// GetAllReviewers godoc
// @Summary Get All Reviewers
// @Description Admin gets list of all activated reviewers
// @Tags Admin - Reviewer Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.APIResponse{data=[]response.ReviewerResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/reviewers [get]
func (ctrl *ReviewerController) GetAllReviewers(c *fiber.Ctx) error {
	result, err := ctrl.service.GetAllReviewers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"Failed to get reviewers",
			err.Error(),
		))
	}

	return c.JSON(response.SuccessResponse(
		"Reviewers retrieved successfully",
		result,
	))
}

// GetAvailablePegawai godoc
// @Summary Get Available Pegawai
// @Description Admin gets list of pegawai with email_umm that can be activated as reviewers
// @Tags Admin - Reviewer Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.APIResponse{data=[]response.AvailablePegawaiResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/reviewers/available [get]
func (ctrl *ReviewerController) GetAvailablePegawai(c *fiber.Ctx) error {
	result, err := ctrl.service.GetAvailablePegawai()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"Failed to get available pegawai",
			err.Error(),
		))
	}

	return c.JSON(response.SuccessResponse(
		"Available pegawai retrieved successfully",
		result,
	))
}

// ActivateReviewer godoc
// @Summary Activate Pegawai as Reviewer
// @Description Admin activates a pegawai to become reviewer
// @Tags Admin - Reviewer Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body request.ActivateReviewerRequest true "Pegawai ID"
// @Success 201 {object} response.APIResponse{data=response.ReviewerResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/reviewers [post]
func (ctrl *ReviewerController) ActivateReviewer(c *fiber.Ctx) error {
	// 1. Parse request body
	var req request.ActivateReviewerRequest
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
	result, err := ctrl.service.ActivateReviewer(&req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to activate reviewer",
			err.Error(),
		))
	}

	// 5. Return success
	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Reviewer berhasil diaktifkan",
		result,
	))
}

// UpdateReviewer godoc
// @Summary Update Reviewer Status
// @Description Admin updates reviewer active status
// @Tags Admin - Reviewer Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Reviewer ID"
// @Param body body request.UpdateReviewerRequest true "Update data"
// @Success 200 {object} response.APIResponse{data=response.ReviewerResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/reviewers/{id} [put]
func (ctrl *ReviewerController) UpdateReviewer(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid reviewer ID",
			err.Error(),
		))
	}

	// 2. Parse request body
	var req request.UpdateReviewerRequest
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
	result, err := ctrl.service.UpdateReviewer(id, &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to update reviewer",
			err.Error(),
		))
	}

	// 6. Return success
	return c.JSON(response.SuccessResponse(
		"Reviewer berhasil diupdate",
		result,
	))
}

// DeleteReviewer godoc
// @Summary Delete Reviewer
// @Description Admin deactivates/deletes reviewer
// @Tags Admin - Reviewer Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Reviewer ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/reviewers/{id} [delete]
func (ctrl *ReviewerController) DeleteReviewer(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid reviewer ID",
			err.Error(),
		))
	}

	// 2. Get user ID
	userID := int(utils.GetCurrentUserID(c))

	// 3. Call service
	if err := ctrl.service.DeleteReviewer(id, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to delete reviewer",
			err.Error(),
		))
	}

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Reviewer berhasil dihapus",
		nil,
	))
}