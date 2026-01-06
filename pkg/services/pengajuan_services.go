package services

import (
	"errors"
	"fmt"
	"time"

	"rires-be/internal/dto/request"
	"rires-be/internal/dto/response"
	"rires-be/internal/models"
	"rires-be/internal/models/external"
	"rires-be/pkg/database"
	"rires-be/pkg/utils"

	"gorm.io/gorm"
)

// PengajuanService handles PKM submission business logic
type PengajuanService struct {
	externalService *ExternalDataService
	fileService     *FileUploadService
	validator       *utils.StatusValidator
	mapper          *MapperService
}

// NewPengajuanService creates a new pengajuan service
func NewPengajuanService() *PengajuanService {
	return &PengajuanService{
		externalService: NewExternalDataService(),
		fileService:     NewFileUploadService(),
		validator:       utils.NewStatusValidator(),
		mapper:          NewMapperService(),
	}
}

// ========================================
// CREATE JUDUL PKM
// ========================================

// CreateJudulPKM creates new PKM title submission
func (s *PengajuanService) CreateJudulPKM(req *request.CreatePengajuanRequest, nimKetua string) (*response.PengajuanResponse, error) {
	// 1. Validate custom rules
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. Check if registration period is open
	if err := s.validator.CanSubmitPengajuan(); err != nil {
		return nil, err
	}

	// 3. Validate team size
	if err := s.validator.ValidateTeamSize(s.convertToAnggotaModels(req.Anggota)); err != nil {
		return nil, err
	}

	// 4. Validate team structure (1 ketua)
	if err := s.validator.ValidateTeamStructure(s.convertToAnggotaModels(req.Anggota)); err != nil {
		return nil, err
	}

	// 5. Validate no duplicate NIM
	if err := s.validator.ValidateNoDuplicateNIM(s.convertToAnggotaModels(req.Anggota)); err != nil {
		return nil, err
	}

	// 6. Find ketua in anggota list
	var ketuaNIM string
	for _, anggota := range req.Anggota {
		if anggota.IsKetua == 1 {
			ketuaNIM = anggota.NIM
			break
		}
	}

	// 7. Verify ketua NIM matches authenticated user
	if ketuaNIM != nimKetua {
		return nil, errors.New("hanya ketua yang dapat membuat pengajuan")
	}

	// 8. Validate all NIMs exist in NEOMAA
	nims := make([]string, len(req.Anggota))
	for i, anggota := range req.Anggota {
		nims[i] = anggota.NIM
	}

	mahasiswaList, err := s.externalService.GetMahasiswaByNIMs(nims)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mahasiswa data: %w", err)
	}

	if len(mahasiswaList) != len(nims) {
		return nil, errors.New("beberapa NIM tidak ditemukan di database mahasiswa")
	}

	// 9. Get kategori
	var kategori models.KategoriPKM
	if err := database.DB.Where("id = ? AND hapus = ?", req.IDKategori, 0).First(&kategori).Error; err != nil {
		return nil, errors.New("kategori PKM tidak ditemukan")
	}

	// 10. Generate kode pengajuan
	tahun := time.Now().Year()
	kodePengajuan, err := utils.GenerateKodePengajuan(&kategori, tahun)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kode pengajuan: %w", err)
	}

	// 11. START TRANSACTION
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 12. Create pengajuan
	now := time.Now()
	pengajuan := &models.Pengajuan{
		KodePengajuan: kodePengajuan,
		NIMKetua:      ketuaNIM,
		IDKategori:    req.IDKategori,
		Judul:         req.Judul,
		Tahun:         tahun,
		StatusJudul:   "PENDING",
		StatusFinal:   "DRAFT",
		Status:        1,
		Hapus:         0,
		TglInsert:     &now,
		UserUpdate:    nimKetua,
	}

	if err := tx.Create(pengajuan).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create pengajuan: %w", err)
	}

	// 13. Create anggota tim
	for _, anggota := range req.Anggota {
		anggotaModel := &models.PengajuanAnggota{
			IDPengajuan: pengajuan.ID,
			NIM:         anggota.NIM,
			IsKetua:     anggota.IsKetua,
			Urutan:      anggota.Urutan,
			Status:      1,
			Hapus:       0,
		}

		if err := tx.Create(anggotaModel).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create anggota tim: %w", err)
		}
	}

	// 14. Create parameter answers
	for _, param := range req.Parameter {
		paramModel := &models.ParameterPKM{
			IDPengajuan: pengajuan.ID,
			IDParameter: param.IDParameter,
			Nilai:       param.Nilai,
		}

		if err := tx.Create(paramModel).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create parameter: %w", err)
		}
	}

	// 15. COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 16. Fetch created data for response
	return s.GetPengajuanDetail(pengajuan.ID)
}

// ========================================
// GET PENGAJUAN DETAIL
// ========================================

// GetPengajuanDetail gets full pengajuan detail with all relations
func (s *PengajuanService) GetPengajuanDetail(idPengajuan int) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Get kategori
	var kategori models.KategoriPKM
	database.DB.Where("id = ?", pengajuan.IDKategori).First(&kategori)

	// 3. Get anggota tim
	var anggotaModels []models.PengajuanAnggota
	database.DB.Where("id_pengajuan = ? AND hapus = ?", pengajuan.ID, 0).
		Order("urutan ASC").
		Find(&anggotaModels)

	// 4. Get mahasiswa data from NEOMAA
	nims := make([]string, len(anggotaModels))
	for i, anggota := range anggotaModels {
		nims[i] = anggota.NIM
	}

	mahasiswaList, _ := s.externalService.GetMahasiswaByNIMs(nims)

	// Find ketua
	var ketua *external.Mahasiswa
	for _, mhs := range mahasiswaList {
		if mhs.KodeSiswa == pengajuan.NIMKetua {
			ketua = &mhs
			break
		}
	}

	// 5. Get parameters
	var parameterList []models.ParameterPKM
	database.DB.Preload("ParameterForm").
		Where("id_pengajuan = ?", pengajuan.ID).
		Find(&parameterList)

	// 6. Get reviewers
	var reviewerJudul, reviewerProposal *external.Pegawai
	if pengajuan.IDPegawaiReviewerJudul != nil {
		reviewerJudul, _ = s.externalService.GetPegawaiByID(*pengajuan.IDPegawaiReviewerJudul)
	}
	if pengajuan.IDPegawaiReviewerProposal != nil {
		reviewerProposal, _ = s.externalService.GetPegawaiByID(*pengajuan.IDPegawaiReviewerProposal)
	}

	// 7. Get review history
	var reviewJudulHistory []models.ReviewJudul
	var reviewProposalHistory []models.ReviewProposal
	
	database.DB.Preload("StatusReview").
		Where("id_pengajuan = ?", pengajuan.ID).
		Order("tgl_review DESC").
		Find(&reviewJudulHistory)
	
	database.DB.Preload("StatusReview").
		Where("id_pengajuan = ?", pengajuan.ID).
		Order("tgl_review DESC").
		Find(&reviewProposalHistory)

	// 8. Map to response DTO
	return s.mapper.MapPengajuanToDetailResponse(
		&pengajuan,
		ketua,
		mahasiswaList,
		anggotaModels,
		&kategori,
		parameterList,
		reviewerJudul,
		reviewerProposal,
		reviewJudulHistory,
		reviewProposalHistory,
	), nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// convertToAnggotaModels converts request anggota to models for validation
func (s *PengajuanService) convertToAnggotaModels(anggota []request.AnggotaRequest) []models.PengajuanAnggota {
	result := make([]models.PengajuanAnggota, len(anggota))
	for i, a := range anggota {
		result[i] = models.PengajuanAnggota{
			NIM:     a.NIM,
			IsKetua: a.IsKetua,
			Urutan:  a.Urutan,
		}
	}
	return result
}