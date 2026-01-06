package controllers

import (
	"rires-be/internal/models"
	"rires-be/pkg/database"
	"rires-be/pkg/services"
	"rires-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type TestHelperController struct{}

func NewTestHelperController() *TestHelperController {
	return &TestHelperController{}
}

// TestExternalData godoc
// @Summary Test External Data Service
// @Description Test querying mahasiswa and pegawai from external databases
// @Tags Test Helpers
// @Accept json
// @Produce json
// @Param nim query string false "NIM Mahasiswa to test"
// @Param id_pegawai query int false "ID Pegawai to test"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Router /test/external-data [get]
func (ctrl *TestHelperController) TestExternalData(c *fiber.Ctx) error {
	externalService := services.NewExternalDataService()
	result := make(map[string]interface{})

	// Test Mahasiswa
	nimTest := c.Query("nim", "202110370311503") // Default test NIM
	mahasiswa, err := externalService.GetMahasiswaByNIM(nimTest)
	if err != nil {
		result["mahasiswa_error"] = err.Error()
		result["mahasiswa"] = nil
	} else {
		result["mahasiswa"] = mahasiswa
	}

	// Test validate NIM
	result["nim_exists"] = externalService.ValidateNIMExists(nimTest)

	// Test Pegawai
	idPegawaiTest := c.QueryInt("id_pegawai", 1) // Default test ID
	pegawai, err := externalService.GetPegawaiByID(idPegawaiTest)
	if err != nil {
		result["pegawai_error"] = err.Error()
		result["pegawai"] = nil
	} else {
		result["pegawai"] = pegawai
	}

	// Test validate Pegawai
	result["pegawai_exists"] = externalService.ValidatePegawaiExists(idPegawaiTest)

	// Test get all reviewers
	reviewers, err := externalService.GetAllReviewers()
	if err != nil {
		result["reviewers_error"] = err.Error()
		result["reviewers_count"] = 0
	} else {
		result["reviewers_count"] = len(reviewers)
		if len(reviewers) > 0 {
			result["reviewers_sample"] = reviewers[0:min(3, len(reviewers))] // First 3
		}
	}

	// Test Prodi & Fakultas
	allProdi, err := externalService.GetAllProdi()
	if err != nil {
		result["prodi_error"] = err.Error()
		result["prodi_count"] = 0
	} else {
		result["prodi_count"] = len(allProdi)
	}

	allFakultas, err := externalService.GetAllFakultas()
	if err != nil {
		result["fakultas_error"] = err.Error()
		result["fakultas_count"] = 0
	} else {
		result["fakultas_count"] = len(allFakultas)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "External Data Service Test",
		"data":    result,
	})
}

// TestCodeGenerator godoc
// @Summary Test Code Generator
// @Description Test generating kode pengajuan
// @Tags Test Helpers
// @Accept json
// @Produce json
// @Param id_kategori query int false "ID Kategori PKM" default(1)
// @Param tahun query int false "Tahun" default(2026)
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Router /test/code-generator [get]
func (ctrl *TestHelperController) TestCodeGenerator(c *fiber.Ctx) error {
	idKategori := c.QueryInt("id_kategori", 1)
	tahun := c.QueryInt("tahun", 2026)

	// Get kategori
	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", idKategori, 0).First(&kategori).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Kategori not found",
		})
	}

	// Generate code
	code, err := utils.GenerateKodePengajuan(&kategori, tahun)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Validate uniqueness
	isUnique := utils.ValidateKodePengajuan(code)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Code Generator Test",
		"data": fiber.Map{
			"kategori":      kategori.NamaKategori,
			"tahun":         tahun,
			"generated_code": code,
			"is_unique":     isUnique,
		},
	})
}

// TestStatusValidator godoc
// @Summary Test Status Validator
// @Description Test status flow validation
// @Tags Test Helpers
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Router /test/status-validator [get]
func (ctrl *TestHelperController) TestStatusValidator(c *fiber.Ctx) error {
	validator := utils.NewStatusValidator()
	result := make(map[string]interface{})

	// Test registration period
	err := validator.CanSubmitPengajuan()
	if err != nil {
		result["can_submit"] = false
		result["submit_error"] = err.Error()
	} else {
		result["can_submit"] = true
	}

	// Test team validation
	testTeam := []models.PengajuanAnggota{
		{NIM: "202110370311503", IsKetua: 1, Urutan: 1},
		{NIM: "202110370311504", IsKetua: 0, Urutan: 2},
		{NIM: "202110370311505", IsKetua: 0, Urutan: 3},
	}

	// Validate team size
	if err := validator.ValidateTeamSize(testTeam); err != nil {
		result["team_size_valid"] = false
		result["team_size_error"] = err.Error()
	} else {
		result["team_size_valid"] = true
	}

	// Validate team structure
	if err := validator.ValidateTeamStructure(testTeam); err != nil {
		result["team_structure_valid"] = false
		result["team_structure_error"] = err.Error()
	} else {
		result["team_structure_valid"] = true
	}

	// Validate no duplicate NIM
	if err := validator.ValidateNoDuplicateNIM(testTeam); err != nil {
		result["no_duplicate_valid"] = false
		result["no_duplicate_error"] = err.Error()
	} else {
		result["no_duplicate_valid"] = true
	}

	// Test with dummy pengajuan
	dummyPengajuan := &models.Pengajuan{
		StatusJudul:    "ACC",
		StatusProposal: "REVISI",
		FileProposal:   "proposal_test.pdf",
		NIMKetua:       "202110370311503",
	}

	result["can_upload_proposal"] = validator.CanSubmitProposal(dummyPengajuan) == nil
	result["can_revise_judul"] = validator.CanReviseJudul(dummyPengajuan) != nil // Should be error (status = ACC, not REVISI)
	result["can_revise_proposal"] = validator.CanReviseProposal(dummyPengajuan) == nil
	result["is_owner"] = validator.IsOwner(dummyPengajuan, "202110370311503")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Status Validator Test",
		"data":    result,
	})
}

// TestFileUpload godoc
// @Summary Test File Upload Service
// @Description Test file upload validation
// @Tags Test Helpers
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Proposal file (PDF/DOC/DOCX, max 2.5MB)"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Router /test/file-upload [post]
func (ctrl *TestHelperController) TestFileUpload(c *fiber.Ctx) error {
	fileUploadService := services.NewFileUploadService()

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
		})
	}

	// Validate file
	if err := fileUploadService.ValidateFile(file); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Upload file (test mode - use dummy kode)
	kodePengajuan := "PKM-TEST-2026-001"
	filename, err := fileUploadService.UploadProposal(file, kodePengajuan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Check if file exists
	exists := fileUploadService.FileExists(filename)

	// Get file size
	size, _ := fileUploadService.GetFileSize(filename)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"original_filename": file.Filename,
			"saved_filename":    filename,
			"file_exists":       exists,
			"file_size":         size,
			"file_size_mb":      float64(size) / (1024 * 1024),
		},
	})
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}