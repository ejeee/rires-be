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

// PengajuanReviewerController handles reviewer PKM review endpoints
type PengajuanReviewerController struct {
	service   *services.PengajuanService
	validator *validator.Validate
}

// NewPengajuanReviewerController creates a new controller instance
func NewPengajuanReviewerController() *PengajuanReviewerController {
	return &PengajuanReviewerController{
		service:   services.NewPengajuanService(),
		validator: validator.New(),
	}
}

// GetMyAssignments godoc
// @Summary Get My Assignments (Reviewer)
// @Description Reviewer gets all pengajuan assigned to them
// @Tags Reviewer - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tipe query string false "Filter by tipe (JUDUL/PROPOSAL/all)" default(all)
// @Success 200 {object} response.APIResponse{data=[]response.PengajuanListResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /reviewer/my-assignments [get]
func (ctrl *PengajuanReviewerController) GetMyAssignments(c *fiber.Ctx) error {
	// 1. Get authenticated reviewer (pegawai ID from JWT)
	// For pegawai, user_id in JWT is their ID from db_user
	// We need to map this to pegawai.id in SIMPEG database
	userID := int(utils.GetCurrentUserID(c))

	// 2. Get filter
	tipeFilter := c.Query("tipe", "all") // JUDUL, PROPOSAL, or all

	// 3. Call service
	result, err := ctrl.service.GetMyAssignments(userID, tipeFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"Failed to get assignments",
			err.Error(),
		))
	}

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Assignments retrieved successfully",
		result,
	))
}

// ReviewJudul godoc
// @Summary Review PKM Title
// @Description Reviewer submits review for PKM title (ACC/REVISI/TOLAK)
// @Tags Reviewer - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param body body request.ReviewJudulRequest true "Review data"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /reviewer/judul/{id}/review [post]
func (ctrl *PengajuanReviewerController) ReviewJudul(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Parse request body
	var req request.ReviewJudulRequest
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

	// 4. Get authenticated reviewer
	userID := int(utils.GetCurrentUserID(c))

	// 5. Call service
	result, err := ctrl.service.ReviewJudul(id, &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to submit review",
			err.Error(),
		))
	}

	// 6. Return success
	return c.JSON(response.SuccessResponse(
		"Review judul berhasil disimpan",
		result,
	))
}

// ReviewProposal godoc
// @Summary Review PKM Proposal
// @Description Reviewer submits review for PKM proposal (ACC/REVISI/TOLAK)
// @Tags Reviewer - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param body body request.ReviewProposalRequest true "Review data"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /reviewer/proposal/{id}/review [post]
func (ctrl *PengajuanReviewerController) ReviewProposal(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Parse request body
	var req request.ReviewProposalRequest
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

	// 4. Get authenticated reviewer
	userID := int(utils.GetCurrentUserID(c))

	// 5. Call service
	result, err := ctrl.service.ReviewProposal(id, &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to submit review",
			err.Error(),
		))
	}

	// 6. Return success
	return c.JSON(response.SuccessResponse(
		"Review proposal berhasil disimpan",
		result,
	))
}

// GetPengajuanDetail godoc
// @Summary Get Pengajuan Detail (Reviewer)
// @Description Reviewer gets detail of pengajuan assigned to them
// @Tags Reviewer - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /reviewer/pengajuan/{id} [get]
func (ctrl *PengajuanReviewerController) GetPengajuanDetail(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Get detail from service
	result, err := ctrl.service.GetPengajuanDetail(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(
			"Pengajuan not found",
			err.Error(),
		))
	}

	// 3. Verify reviewer has access (assigned to them)
	userID := int(utils.GetCurrentUserID(c))
	
	// Check if reviewer is assigned to this pengajuan
	// (This validation will be done in service layer later if needed)
	_ = userID // TODO: Implement access check

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Pengajuan detail",
		result,
	))
}