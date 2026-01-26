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
	// 1. Get authenticated reviewer (id_pegawai from JWT user_data)
	userData := utils.GetCurrentUserData(c)
	idPegawai, _ := strconv.Atoi(userData["id_pegawai"])

	if idPegawai == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(
			"Reviewer ID not found in token. Please relogin.",
			"",
		))
	}

	// 2. Get filter
	tipeFilter := c.Query("tipe", "all") // JUDUL, PROPOSAL, or all

	// 3. Call service
	result, err := ctrl.service.GetMyAssignments(idPegawai, tipeFilter)
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

	// 4. Get authenticated reviewer and check if admin
	userData := utils.GetCurrentUserData(c)
	idPegawai, _ := strconv.Atoi(userData["id_pegawai"])
	isAdmin := utils.IsAdmin(c)

	if idPegawai == 0 && !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(
			"Reviewer ID not found in token. Please relogin.",
			"",
		))
	}

	// 5. Call service
	result, err := ctrl.service.ReviewJudul(id, &req, idPegawai, isAdmin)
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

	// 4. Get authenticated reviewer and check if admin
	userData := utils.GetCurrentUserData(c)
	idPegawai, _ := strconv.Atoi(userData["id_pegawai"])
	isAdmin := utils.IsAdmin(c)

	if idPegawai == 0 && !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(
			"Reviewer ID not found in token. Please relogin.",
			"",
		))
	}

	// 5. Call service
	result, err := ctrl.service.ReviewProposal(id, &req, idPegawai, isAdmin)
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
	userData := utils.GetCurrentUserData(c)
	idPegawai, _ := strconv.Atoi(userData["id_pegawai"])
	isAdmin := utils.IsAdmin(c)

	// Check if reviewer is assigned to this pengajuan
	if !isAdmin {
		isAssigned := false
		if result.ReviewerJudul != nil && result.ReviewerJudul.ID == idPegawai {
			isAssigned = true
		}
		if result.ReviewerProposal != nil && result.ReviewerProposal.ID == idPegawai {
			isAssigned = true
		}

		if !isAssigned {
			return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse(
				"Access denied. You are not assigned to this pengajuan.",
				"",
			))
		}
	}

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Pengajuan detail",
		result,
	))
}

// CancelReviewJudul godoc
// @Summary Cancel Review PKM Title
// @Description Cancel/reset review for PKM title (status back to ON_REVIEW)
// @Tags Reviewer - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /reviewer/judul/{id}/cancel-review [post]
func (ctrl *PengajuanReviewerController) CancelReviewJudul(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Get authenticated user and check if admin
	userData := utils.GetCurrentUserData(c)
	idPegawai, _ := strconv.Atoi(userData["id_pegawai"])
	isAdmin := utils.IsAdmin(c)

	if idPegawai == 0 && !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(
			"Reviewer ID not found in token. Please relogin.",
			"",
		))
	}

	// 3. Call service
	result, err := ctrl.service.CancelReviewJudul(id, idPegawai, isAdmin)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to cancel review",
			err.Error(),
		))
	}

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Review judul berhasil dibatalkan",
		result,
	))
}

// CancelReviewProposal godoc
// @Summary Cancel Review Proposal
// @Description Reviewer/admin cancels submitted review for proposal (resets back to ON_REVIEW)
// @Tags Reviewer - Review Proposal
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /reviewer/proposal/{id}/cancel-review [post]
func (ctrl *PengajuanReviewerController) CancelReviewProposal(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Get authenticated user and check if admin
	userData := utils.GetCurrentUserData(c)
	idPegawai, _ := strconv.Atoi(userData["id_pegawai"])
	isAdmin := utils.IsAdmin(c)

	if idPegawai == 0 && !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(
			"Reviewer ID not found in token. Please relogin.",
			"",
		))
	}

	// 3. Call service
	result, err := ctrl.service.CancelReviewProposal(id, idPegawai, isAdmin)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to cancel review",
			err.Error(),
		))
	}

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Review proposal berhasil dibatalkan",
		result,
	))
}
