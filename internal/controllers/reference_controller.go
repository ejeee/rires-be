package controllers

import (
	"rires-be/internal/dto/response"
	"rires-be/internal/models/external"
	"rires-be/pkg/database"
	"rires-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ReferenceController struct{}

func NewReferenceController() *ReferenceController {
	return &ReferenceController{}
}

// GetAllFakultas godoc
// @Summary Get All Fakultas
// @Description Get list of all fakultas from NEOMAAREF database
// @Tags Reference Data
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string,data=response.FakultasListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /reference/fakultas [get]
func (ctrl *ReferenceController) GetAllFakultas(c *fiber.Ctx) error {
	var fakultasList []external.Fakultas

	// Query from NEOMAAREF database
	if err := database.DBNeomaaRef.Where("hapus = ? AND st_aktif = ?", 0, 1).
		Order("namaFakultas ASC").
		Find(&fakultasList).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch fakultas data")
	}

	// Transform to response
	var data []response.FakultasRefResponse
	for _, fak := range fakultasList {
		data = append(data, response.FakultasRefResponse{
			Kode:          fak.Kode,
			NamaFakultas:  fak.NamaFakultas,
			NamaFakPendek: fak.NamaFakPendek,
		})
	}

	result := response.FakultasListResponse{
		Data:  data,
		Total: len(data),
	}

	return utils.SuccessResponse(c, "Fakultas retrieved successfully", result)
}

// GetAllProdi godoc
// @Summary Get All Program Studi
// @Description Get list of all program studi from NEOMAAREF database
// @Tags Reference Data
// @Accept json
// @Produce json
// @Param fakultas query string false "Filter by fakultas code"
// @Success 200 {object} object{success=bool,message=string,data=response.ProdiListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /reference/prodi [get]
func (ctrl *ReferenceController) GetAllProdi(c *fiber.Ctx) error {
	fakultasKode := c.Query("fakultas", "")

	var prodiList []external.Prodi

	// Build query
	query := database.DBNeomaaRef.Where("hapus = ?", 0)

	// Filter by fakultas if provided
	if fakultasKode != "" {
		query = query.Where("kodeFakultas = ?", fakultasKode)
	}

	// Query from NEOMAAREF database with Fakultas preload
	if err := query.Preload("Fakultas").
		Order("nama_depart ASC").
		Find(&prodiList).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch prodi data")
	}

	// Transform to response
	var data []response.ProdiRefResponse
	for _, prodi := range prodiList {
		prodiResp := response.ProdiRefResponse{
			Kode:         prodi.Kode,
			KodeFakultas: prodi.KodeFakultas,
			KodeDepart:   prodi.KodeDepart,
			NamaDepart:   prodi.NamaDepart,
			NamaSingkat:  prodi.NamaSingkat,
		}

		// Add fakultas if available
		if prodi.Fakultas != nil {
			prodiResp.Fakultas = &response.FakultasRefResponse{
				Kode:          prodi.Fakultas.Kode,
				NamaFakultas:  prodi.Fakultas.NamaFakultas,
				NamaFakPendek: prodi.Fakultas.NamaFakPendek,
			}
		}

		data = append(data, prodiResp)
	}

	result := response.ProdiListResponse{
		Data:  data,
		Total: len(data),
	}

	return utils.SuccessResponse(c, "Program studi retrieved successfully", result)
}

// GetProdiByFakultas godoc
// @Summary Get Program Studi by Fakultas
// @Description Get list of program studi by fakultas code from NEOMAAREF database
// @Tags Reference Data
// @Accept json
// @Produce json
// @Param kode path string true "Fakultas Code"
// @Success 200 {object} object{success=bool,message=string,data=response.ProdiListResponse}
// @Failure 500 {object} object{success=bool,message=string}
// @Security BearerAuth
// @Router /reference/prodi/fakultas/{kode} [get]
func (ctrl *ReferenceController) GetProdiByFakultas(c *fiber.Ctx) error {
	kode := c.Params("kode")

	var prodiList []external.Prodi

	// Query from NEOMAAREF database
	if err := database.DBNeomaaRef.Where("hapus = ? AND kodeFakultas = ?", 0, kode).
		Preload("Fakultas").
		Order("nama_depart ASC").
		Find(&prodiList).Error; err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch prodi data")
	}

	// Transform to response
	var data []response.ProdiRefResponse
	for _, prodi := range prodiList {
		prodiResp := response.ProdiRefResponse{
			Kode:         prodi.Kode,
			KodeFakultas: prodi.KodeFakultas,
			KodeDepart:   prodi.KodeDepart,
			NamaDepart:   prodi.NamaDepart,
			NamaSingkat:  prodi.NamaSingkat,
		}

		// Add fakultas if available
		if prodi.Fakultas != nil {
			prodiResp.Fakultas = &response.FakultasRefResponse{
				Kode:          prodi.Fakultas.Kode,
				NamaFakultas:  prodi.Fakultas.NamaFakultas,
				NamaFakPendek: prodi.Fakultas.NamaFakPendek,
			}
		}

		data = append(data, prodiResp)
	}

	result := response.ProdiListResponse{
		Data:  data,
		Total: len(data),
	}

	return utils.SuccessResponse(c, "Program studi retrieved successfully", result)
}
