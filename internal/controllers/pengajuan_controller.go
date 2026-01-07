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

// PengajuanController handles mahasiswa PKM submission endpoints
type PengajuanController struct {
	service   *services.PengajuanService
	validator *validator.Validate
}

// NewPengajuanController creates a new controller instance
func NewPengajuanController() *PengajuanController {
	return &PengajuanController{
		service:   services.NewPengajuanService(),
		validator: validator.New(),
	}
}

// CreateJudulPKM godoc
// @Summary Create PKM Title Submission
// @Description Mahasiswa (ketua) creates new PKM title submission with team members
// @Tags Mahasiswa - Pengajuan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body request.CreatePengajuanRequest true "Pengajuan data"
// @Success 201 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /pengajuan/judul [post]
func (ctrl *PengajuanController) CreateJudulPKM(c *fiber.Ctx) error {
	// 1. Parse request body
	var req request.CreatePengajuanRequest
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
			ctrl.formatValidationErrors(err),
		))
	}

	// 3. Custom validation
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Validation failed",
			err.Error(),
		))
	}

	// 4. Get authenticated user (NIM ketua)
	nimKetua := utils.GetCurrentUsername(c) // Assuming JWT contains NIM in username field

	// 5. Call service
	result, err := ctrl.service.CreateJudulPKM(&req, nimKetua)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"Failed to create pengajuan",
			err.Error(),
		))
	}

	// 6. Return success response
	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Pengajuan berhasil dibuat",
		result,
	))
}

// GetMySubmissions godoc
// @Summary Get My Submissions
// @Description Get all submissions created by authenticated mahasiswa
// @Tags Mahasiswa - Pengajuan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param status query string false "Filter by status (all/draft/pending/acc/revisi/tolak)"
// @Success 200 {object} response.APIResponse{data=[]response.PengajuanListResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /pengajuan/my-submissions [get]
func (ctrl *PengajuanController) GetMySubmissions(c *fiber.Ctx) error {
	// TODO: Implement list my submissions
	// Get NIM from JWT
	// Filter by NIM ketua
	// Return list

	return c.JSON(response.SuccessResponse(
		"Feature coming soon",
		nil,
	))
}

// GetPengajuanDetail godoc
// @Summary Get Pengajuan Detail
// @Description Get detailed information of a pengajuan (only if user is part of the team)
// @Tags Mahasiswa - Pengajuan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /pengajuan/{id} [get]
func (ctrl *PengajuanController) GetPengajuanDetail(c *fiber.Ctx) error {
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

	// 3. Check authorization - only team members can view
	// TODO: Verify user is part of the team
	// nimUser := utils.GetCurrentUsername(c)
	// if !isTeamMember(result, nimUser) {
	//     return c.Status(403).JSON(...)
	// }

	// 4. Return success
	return c.JSON(response.SuccessResponse(
		"Pengajuan detail",
		result,
	))
}

// UpdateJudul godoc
// @Summary Update/Revise PKM Title
// @Description Ketua revises PKM title (only allowed when status = REVISI)
// @Tags Mahasiswa - Pengajuan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param body body request.UpdateJudulRequest true "Updated title data"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /pengajuan/judul/{id} [put]
func (ctrl *PengajuanController) UpdateJudul(c *fiber.Ctx) error {
	// TODO: Implement update judul
	// 1. Parse ID & request body
	// 2. Validate request
	// 3. Check if status = REVISI
	// 4. Check if user is ketua
	// 5. Update judul & parameter
	// 6. Return updated detail

	return c.JSON(response.SuccessResponse(
		"Feature coming soon",
		nil,
	))
}

// UploadProposal godoc
// @Summary Upload Proposal File
// @Description Ketua uploads proposal PDF/DOC/DOCX (only allowed when status_judul = ACC)
// @Tags Mahasiswa - Pengajuan
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param file formData file true "Proposal file (PDF/DOC/DOCX, max 2.5MB)"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /pengajuan/{id}/proposal [post]
func (ctrl *PengajuanController) UploadProposal(c *fiber.Ctx) error {
	// TODO: Implement upload proposal
	// 1. Parse ID
	// 2. Get file from form
	// 3. Validate file (extension, size)
	// 4. Check if status_judul = ACC
	// 5. Check if user is ketua
	// 6. Upload file using FileUploadService
	// 7. Update pengajuan.file_proposal
	// 8. Set status_proposal = PENDING
	// 9. Return updated detail

	return c.JSON(response.SuccessResponse(
		"Feature coming soon",
		nil,
	))
}

// ReviseProposal godoc
// @Summary Revise Proposal File
// @Description Ketua revises proposal by uploading new file (only allowed when status_proposal = REVISI)
// @Tags Mahasiswa - Pengajuan
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Pengajuan ID"
// @Param file formData file true "Revised proposal file (PDF/DOC/DOCX, max 2.5MB)"
// @Success 200 {object} response.APIResponse{data=response.PengajuanResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security BearerAuth
// @Router /pengajuan/{id}/proposal [put]
func (ctrl *PengajuanController) ReviseProposal(c *fiber.Ctx) error {
	// TODO: Implement revise proposal
	// Similar to UploadProposal but:
	// 1. Check if status_proposal = REVISI
	// 2. Delete old file
	// 3. Upload new file
	// 4. Update file_proposal
	// 5. Reset status_proposal = PENDING

	return c.JSON(response.SuccessResponse(
		"Feature coming soon",
		nil,
	))
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// formatValidationErrors formats validator errors to readable format
func (ctrl *PengajuanController) formatValidationErrors(err error) []response.ValidationErrorResponse {
	var errors []response.ValidationErrorResponse

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, response.ValidationErrorResponse{
				Field:   e.Field(),
				Message: ctrl.getValidationMessage(e),
			})
		}
	}

	return errors
}

// getValidationMessage returns user-friendly validation message
func (ctrl *PengajuanController) getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "min":
		return e.Field() + " must be at least " + e.Param()
	case "max":
		return e.Field() + " must be at most " + e.Param()
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}