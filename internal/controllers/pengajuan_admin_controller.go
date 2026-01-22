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

// PengajuanAdminController handles admin PKM submission management
type PengajuanAdminController struct {
	service   *services.PengajuanService
	validator *validator.Validate
}

// NewPengajuanAdminController creates a new controller instance
func NewPengajuanAdminController() *PengajuanAdminController {
	return &PengajuanAdminController{
		service:   services.NewPengajuanService(),
		validator: validator.New(),
	}
}

// GetAllPengajuan godoc
// @Summary Get All Pengajuan (Admin)
// @Description Admin gets all pengajuan with filters and pagination
// @Tags Admin - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param status_judul query string false "Filter by status judul"
// @Param status_proposal query string false "Filter by status proposal"
// @Param status_final query string false "Filter by status final"
// @Param id_kategori query int false "Filter by kategori"
// @Param tahun query int false "Filter by tahun"
// @Success 200 {object} response.APIResponse{data=response.PaginatedResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/pengajuan [get]
func (ctrl *PengajuanAdminController) GetAllPengajuan(c *fiber.Ctx) error {
	// 1. Parse query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	statusJudul := c.Query("status_judul", "")
	statusProposal := c.Query("status_proposal", "")
	statusFinal := c.Query("status_final", "")
	idKategori, _ := strconv.Atoi(c.Query("id_kategori", "0"))
	tahun, _ := strconv.Atoi(c.Query("tahun", "0"))

	// 2. Build filters
	filters := map[string]interface{}{
		"page":            page,
		"per_page":        perPage,
		"status_judul":    statusJudul,
		"status_proposal": statusProposal,
		"status_final":    statusFinal,
		"id_kategori":     idKategori,
		"tahun":           tahun,
	}

	// 3. Call service
	result, pagination, err := ctrl.service.GetAllPengajuan(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"Failed to get pengajuan list",
			err.Error(),
		))
	}

	// 4. Return paginated response
	return c.JSON(response.SuccessResponse(
		"Pengajuan list retrieved successfully",
		response.PaginatedResponse{
			Data:       result,
			Pagination: pagination,
		},
	))
}

// GetPengajuanDetail godoc
// @Summary Get Pengajuan Detail (Admin)
// @Description Admin gets full detail of any pengajuan
// @Tags Admin - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/pengajuan/{id} [get]
func (ctrl *PengajuanAdminController) GetPengajuanDetail(c *fiber.Ctx) error {
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

	// 3. Return success
	return c.JSON(response.SuccessResponse(
		"Pengajuan detail",
		result,
	))
}

// AssignReviewerJudul godoc
// @Summary Assign Reviewer for Judul
// @Description Admin assigns reviewer (pegawai) to review PKM title
// @Tags Admin - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param body body request.AssignReviewerRequest true "Reviewer data"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/pengajuan/{id}/assign-reviewer-judul [post]
func (ctrl *PengajuanAdminController) AssignReviewerJudul(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Parse request body
	var req request.AssignReviewerRequest
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

	// 4. Get user ID for audit
	userID := int(utils.GetCurrentUserID(c))

	// 5. Call service
	result, err := ctrl.service.AssignReviewerJudul(id, req.IDReviewer, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to assign reviewer",
			err.Error(),
		))
	}

	// 6. Return success
	return c.JSON(response.SuccessResponse(
		"Reviewer berhasil di-assign untuk review judul",
		result,
	))
}

// AssignReviewerProposal godoc
// @Summary Assign Reviewer for Proposal
// @Description Admin assigns reviewer (pegawai) to review PKM proposal
// @Tags Admin - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param body body request.AssignReviewerRequest true "Reviewer data"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/pengajuan/{id}/assign-reviewer-proposal [post]
func (ctrl *PengajuanAdminController) AssignReviewerProposal(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Parse request body
	var req request.AssignReviewerRequest
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

	// 4. Get user ID for audit
	userID := int(utils.GetCurrentUserID(c))

	// 5. Call service
	result, err := ctrl.service.AssignReviewerProposal(id, req.IDReviewer, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to assign reviewer",
			err.Error(),
		))
	}

	// 6. Return success
	return c.JSON(response.SuccessResponse(
		"Reviewer berhasil di-assign untuk review proposal",
		result,
	))
}

// AnnounceFinalResult godoc
// @Summary Announce Final Result
// @Description Admin announces final result (LOLOS/TIDAK_LOLOS)
// @Tags Admin - Pengajuan PKM
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param body body request.AnnounceRequest true "Final result"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /admin/pengajuan/{id}/announce [post]
func (ctrl *PengajuanAdminController) AnnounceFinalResult(c *fiber.Ctx) error {
	// 1. Parse ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Invalid pengajuan ID",
			err.Error(),
		))
	}

	// 2. Parse request body
	var req request.AnnounceRequest
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

	// 4. Get user ID for audit
	userID := int(utils.GetCurrentUserID(c))

	// 5. Call service
	result, err := ctrl.service.AnnounceFinalResult(id, req.StatusFinal, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to announce result",
			err.Error(),
		))
	}

	// 6. Return success
	return c.JSON(response.SuccessResponse(
		"Pengumuman final berhasil",
		result,
	))
}
