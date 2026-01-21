package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
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
func (s *PengajuanService) CreateJudulPKM(req *request.CreatePengajuanRequest, nimKetua string, isAdmin bool) (*response.PengajuanResponse, error) {
	// 1. Auto-add ketua to anggota list if not already present
	ketuaFound := false
	for i, anggota := range req.Anggota {
		if anggota.NIM == nimKetua {
			// Mark this anggota as ketua
			req.Anggota[i].IsKetua = 1
			req.Anggota[i].Urutan = 1
			ketuaFound = true
		} else if anggota.IsKetua == 1 {
			// Reset other is_ketua flags
			req.Anggota[i].IsKetua = 0
		}
	}

	// If ketua not in list, add them
	if !ketuaFound {
		ketuaAnggota := request.AnggotaRequest{
			NIM:     nimKetua,
			IsKetua: 1,
			Urutan:  1,
		}
		// Prepend ketua to list
		req.Anggota = append([]request.AnggotaRequest{ketuaAnggota}, req.Anggota...)
	}

	// 2. Fix urutan for anggota (ketua = 1, others = 2, 3, 4...)
	urutanCounter := 2
	for i := range req.Anggota {
		if req.Anggota[i].IsKetua != 1 {
			req.Anggota[i].Urutan = urutanCounter
			urutanCounter++
		}
	}

	// 3. Validate custom rules (now ketua is guaranteed to exist)
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 4. Check if registration period is open (skip for admin)
	if !isAdmin {
		if err := s.validator.CanSubmitPengajuan(); err != nil {
			return nil, err
		}
	}

	// 5. Validate team size
	if err := s.validator.ValidateTeamSize(s.convertToAnggotaModels(req.Anggota)); err != nil {
		return nil, err
	}

	// 6. Validate team structure (1 ketua)
	if err := s.validator.ValidateTeamStructure(s.convertToAnggotaModels(req.Anggota)); err != nil {
		return nil, err
	}

	// 7. Validate no duplicate NIM
	if err := s.validator.ValidateNoDuplicateNIM(s.convertToAnggotaModels(req.Anggota)); err != nil {
		return nil, err
	}

	// 8. Get ketua NIM (now guaranteed to be nimKetua)
	ketuaNIM := nimKetua

	// 9. Validate all NIMs exist in NEOMAA (skip for admin - for testing purposes)
	nims := make([]string, len(req.Anggota))
	for i, anggota := range req.Anggota {
		nims[i] = anggota.NIM
	}

	var mahasiswaList []external.Mahasiswa
	if !isAdmin {
		// For mahasiswa: validate NIMs in NEOMAA
		var err error
		mahasiswaList, err = s.externalService.GetMahasiswaByNIMs(nims)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch mahasiswa data: %w", err)
		}

		if len(mahasiswaList) != len(nims) {
			return nil, errors.New("beberapa NIM tidak ditemukan di database mahasiswa")
		}
	} else {
		// For admin: try to fetch but don't fail if not found
		mahasiswaList, _ = s.externalService.GetMahasiswaByNIMs(nims)
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

	// Convert parameter_data map to JSON string
	var parameterDataJSON string
	if req.ParameterData != nil {
		paramDataBytes, err := json.Marshal(req.ParameterData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameter_data: %w", err)
		}
		parameterDataJSON = string(paramDataBytes)
	}

	pengajuan := &models.Pengajuan{
		KodePengajuan:   kodePengajuan,
		NamaKetua:       req.NamaKetua,
		NIMKetua:        ketuaNIM,
		IDKategori:      req.IDKategori,
		Judul:           req.Judul,
		EmailKetua:      req.EmailKetua,
		NoHPKetua:       req.NoHPKetua,
		ProgramStudi:    req.ProgramStudi,
		Fakultas:        req.Fakultas,
		DosenPembimbing: req.DosenPembimbing,
		ParameterData:   parameterDataJSON,
		TglPengajuan:    &now,
		Tahun:           tahun,
		StatusJudul:     "PENDING",
		StatusFinal:     "DRAFT",
		Status:          1,
		Hapus:           0,
		TglInsert:       &now,
		UserUpdate:      nimKetua, // Store NIM for mahasiswa
	}

	if err := tx.Create(pengajuan).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create pengajuan: %w", err)
	}

	// 13. Create anggota tim
	for _, anggota := range req.Anggota {
		anggotaModel := &models.PengajuanAnggota{
			IDPengajuan: pengajuan.ID,
			NIMAnggota:  anggota.NIM,
			NamaAnggota: anggota.Nama,
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

	// 14. COMMIT TRANSACTION
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
		nims[i] = anggota.NIMAnggota
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
// GET MY SUBMISSIONS
// ========================================

// GetMySubmissions gets all submissions by authenticated mahasiswa (ketua)
func (s *PengajuanService) GetMySubmissions(nimKetua string, statusFilter string) ([]response.PengajuanListResponse, error) {
	// Build query
	query := database.DB.Where("nim_ketua = ? AND hapus = ?", nimKetua, 0)

	// Apply status filter
	switch statusFilter {
	case "pending":
		query = query.Where("status_judul = ?", "PENDING")
	case "acc":
		query = query.Where("status_judul = ?", "ACC")
	case "revisi":
		query = query.Where("status_judul = ?", "REVISI")
	case "tolak":
		query = query.Where("status_judul = ?", "TOLAK")
		// "all" = no filter
	}

	// Get pengajuan list
	var pengajuanList []models.Pengajuan
	if err := query.Order("tgl_insert DESC").Find(&pengajuanList).Error; err != nil {
		return nil, err
	}

	// Build response list
	result := make([]response.PengajuanListResponse, 0)
	for _, pengajuan := range pengajuanList {
		// Get kategori
		var kategori models.KategoriPKM
		database.DB.Where("id = ?", pengajuan.IDKategori).First(&kategori)

		// Get mahasiswa ketua
		ketua, _ := s.externalService.GetMahasiswaByNIM(pengajuan.NIMKetua)

		// Count anggota
		var anggotaCount int64
		database.DB.Model(&models.PengajuanAnggota{}).
			Where("id_pengajuan = ? AND hapus = ?", pengajuan.ID, 0).
			Count(&anggotaCount)

		// Get reviewer proposal (if assigned)
		var reviewerProposal *external.Pegawai
		if pengajuan.IDPegawaiReviewerProposal != nil {
			reviewerProposal, _ = s.externalService.GetPegawaiByID(*pengajuan.IDPegawaiReviewerProposal)
		}

		// Map to list response
		listResp := s.mapper.MapPengajuanToListResponse(
			&pengajuan,
			ketua,
			&kategori,
			int(anggotaCount),
			reviewerProposal,
		)

		result = append(result, *listResp)
	}

	return result, nil
}

// ========================================
// UPDATE JUDUL
// ========================================

// UpdateJudul updates/revises PKM title
func (s *PengajuanService) UpdateJudul(idPengajuan int, req *request.UpdateJudulRequest, nimKetua string) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Check if user is ketua
	if !pengajuan.IsOwner(nimKetua) {
		return nil, errors.New("hanya ketua yang dapat merevisi judul")
	}

	// 3. Check if can revise (status must be REVISI)
	if !pengajuan.CanReviseJudul() {
		return nil, errors.New("judul hanya dapat direvisi jika status = REVISI")
	}

	// 4. START TRANSACTION
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 5. Convert parameter_data map to JSON string
	var parameterDataJSON string
	if req.ParameterData != nil {
		paramDataBytes, err := json.Marshal(req.ParameterData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameter_data: %w", err)
		}
		parameterDataJSON = string(paramDataBytes)
	}

	// 6. Update judul and biodata
	updates := map[string]interface{}{
		"judul":        req.Judul,
		"status_judul": "PENDING", // Reset to PENDING after revision
		"user_update":  nimKetua,
	}

	// Update biodata if provided
	if req.NamaKetua != "" {
		updates["nama_ketua"] = req.NamaKetua
	}
	if req.EmailKetua != "" {
		updates["email_ketua"] = req.EmailKetua
	}
	if req.NoHPKetua != "" {
		updates["no_hp_ketua"] = req.NoHPKetua
	}
	if req.ProgramStudi != "" {
		updates["program_studi"] = req.ProgramStudi
	}
	if req.Fakultas != "" {
		updates["fakultas"] = req.Fakultas
	}
	if req.DosenPembimbing != "" {
		updates["dosen_pembimbing"] = req.DosenPembimbing
	}
	if req.IDKategori != 0 {
		updates["id_kategori"] = req.IDKategori
	}
	if parameterDataJSON != "" {
		updates["parameter_data"] = parameterDataJSON
	}

	if err := tx.Model(&pengajuan).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 7. Update anggota if provided
	if len(req.Anggota) > 0 {
		// Delete old anggota
		if err := tx.Where("id_pengajuan = ?", pengajuan.ID).Delete(&models.PengajuanAnggota{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Create new anggota
		for _, anggota := range req.Anggota {
			anggotaModel := &models.PengajuanAnggota{
				IDPengajuan: pengajuan.ID,
				NIMAnggota:  anggota.NIM,
				NamaAnggota: anggota.Nama,
				IsKetua:     anggota.IsKetua,
				Urutan:      anggota.Urutan,
				Status:      1,
				Hapus:       0,
			}

			if err := tx.Create(anggotaModel).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 8. COMMIT
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 9. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// ========================================
// UPLOAD PROPOSAL
// ========================================

// UploadProposal uploads proposal file
func (s *PengajuanService) UploadProposal(idPengajuan int, file *multipart.FileHeader, nimKetua string) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Check if user is ketua
	if !pengajuan.IsOwner(nimKetua) {
		return nil, errors.New("hanya ketua yang dapat upload proposal")
	}

	// 3. Check if can upload (status_judul must be ACC)
	if !pengajuan.CanUploadProposal() {
		return nil, errors.New("proposal hanya dapat diupload jika judul sudah ACC")
	}

	// 4. Upload file using FileUploadService
	filename, err := s.fileService.UploadProposal(file, pengajuan.KodePengajuan)
	if err != nil {
		return nil, fmt.Errorf("gagal upload file: %w", err)
	}

	// 5. Update pengajuan
	updates := map[string]interface{}{
		"file_proposal":   filename,
		"status_proposal": "PENDING", // Set to PENDING after upload
		"user_update":     nimKetua,
	}

	if err := database.DB.Model(&pengajuan).Updates(updates).Error; err != nil {
		// Delete uploaded file if DB update fails
		s.fileService.DeleteFile(filename)
		return nil, err
	}

	// 6. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// ========================================
// REVISE PROPOSAL
// ========================================

// ReviseProposal revises proposal file
func (s *PengajuanService) ReviseProposal(idPengajuan int, file *multipart.FileHeader, nimKetua string) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Check if user is ketua
	if !pengajuan.IsOwner(nimKetua) {
		return nil, errors.New("hanya ketua yang dapat merevisi proposal")
	}

	// 3. Check if can revise (status_proposal must be REVISI)
	if !pengajuan.CanReviseProposal() {
		return nil, errors.New("proposal hanya dapat direvisi jika status = REVISI")
	}

	// 4. Delete old file
	if pengajuan.FileProposal != "" {
		s.fileService.DeleteFile(pengajuan.FileProposal)
	}

	// 5. Upload new file
	filename, err := s.fileService.UploadProposal(file, pengajuan.KodePengajuan)
	if err != nil {
		return nil, fmt.Errorf("gagal upload file: %w", err)
	}

	// 6. Update pengajuan
	updates := map[string]interface{}{
		"file_proposal":   filename,
		"status_proposal": "PENDING", // Reset to PENDING after revision
		"user_update":     nimKetua,
	}

	if err := database.DB.Model(&pengajuan).Updates(updates).Error; err != nil {
		// Delete uploaded file if DB update fails
		s.fileService.DeleteFile(filename)
		return nil, err
	}

	// 7. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// ========================================
// ADMIN - GET ALL PENGAJUAN
// ========================================

// GetAllPengajuan gets all pengajuan with filters and pagination (admin only)
func (s *PengajuanService) GetAllPengajuan(filters map[string]interface{}) ([]response.PengajuanListResponse, *response.PaginationResponse, error) {
	// 1. Parse filters
	page := filters["page"].(int)
	perPage := filters["per_page"].(int)
	statusJudul := filters["status_judul"].(string)
	statusProposal := filters["status_proposal"].(string)
	statusFinal := filters["status_final"].(string)
	idKategori := filters["id_kategori"].(int)
	tahun := filters["tahun"].(int)

	// 2. Build query
	query := database.DB.Where("hapus = ?", 0)

	// Apply filters
	if statusJudul != "" {
		query = query.Where("status_judul = ?", statusJudul)
	}
	if statusProposal != "" {
		query = query.Where("status_proposal = ?", statusProposal)
	}
	if statusFinal != "" {
		query = query.Where("status_final = ?", statusFinal)
	}
	if idKategori > 0 {
		query = query.Where("id_kategori = ?", idKategori)
	}
	if tahun > 0 {
		query = query.Where("tahun = ?", tahun)
	}

	// 3. Count total records
	var totalRecords int64
	query.Model(&models.Pengajuan{}).Count(&totalRecords)

	// 4. Apply pagination
	offset := (page - 1) * perPage
	query = query.Limit(perPage).Offset(offset).Order("tgl_insert DESC")

	// 5. Get pengajuan list
	var pengajuanList []models.Pengajuan
	if err := query.Find(&pengajuanList).Error; err != nil {
		return nil, nil, err
	}

	// 6. Build response list
	result := make([]response.PengajuanListResponse, 0)
	for _, pengajuan := range pengajuanList {
		// Get kategori
		var kategori models.KategoriPKM
		database.DB.Where("id = ?", pengajuan.IDKategori).First(&kategori)

		// Get mahasiswa ketua
		ketua, _ := s.externalService.GetMahasiswaByNIM(pengajuan.NIMKetua)

		// Count anggota
		var anggotaCount int64
		database.DB.Model(&models.PengajuanAnggota{}).
			Where("id_pengajuan = ? AND hapus = ?", pengajuan.ID, 0).
			Count(&anggotaCount)

		// Get reviewer proposal (if assigned)
		var reviewerProposal *external.Pegawai
		if pengajuan.IDPegawaiReviewerProposal != nil {
			reviewerProposal, _ = s.externalService.GetPegawaiByID(*pengajuan.IDPegawaiReviewerProposal)
		}

		// Map to list response
		listResp := s.mapper.MapPengajuanToListResponse(
			&pengajuan,
			ketua,
			&kategori,
			int(anggotaCount),
			reviewerProposal,
		)

		result = append(result, *listResp)
	}

	// 7. Build pagination response
	paginationResp := response.NewPaginationResponse(page, perPage, totalRecords)

	return result, paginationResp, nil
}

// ========================================
// ADMIN - ASSIGN REVIEWER
// ========================================

// AssignReviewerJudul assigns reviewer for title review
func (s *PengajuanService) AssignReviewerJudul(idPengajuan int, idPegawai int, userID int) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Validate pegawai exists
	if !s.externalService.ValidatePegawaiExists(idPegawai) {
		return nil, errors.New("pegawai/reviewer tidak ditemukan")
	}

	// 3. Check if status allows assignment (must be PENDING or ON_REVIEW)
	if pengajuan.StatusJudul != "PENDING" && pengajuan.StatusJudul != "ON_REVIEW" {
		return nil, errors.New("reviewer hanya dapat di-assign untuk pengajuan dengan status PENDING atau ON_REVIEW")
	}

	// 4. Update pengajuan
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"id_pegawai_reviewer_judul": idPegawai,
		"status_judul":              "ON_REVIEW",
		"user_update":               userUpdateStr,
	}

	if err := database.DB.Model(&pengajuan).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 5. Create plotting record
	plotting := &models.PlottingReviewer{
		IDPengajuan: pengajuan.ID,
		IDPegawai:   idPegawai,
		Tipe:        "JUDUL",
	}

	database.DB.Create(plotting)

	// 6. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// AssignReviewerProposal assigns reviewer for proposal review
func (s *PengajuanService) AssignReviewerProposal(idPengajuan int, idPegawai int, userID int) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Validate pegawai exists
	if !s.externalService.ValidatePegawaiExists(idPegawai) {
		return nil, errors.New("pegawai/reviewer tidak ditemukan")
	}

	// 3. Check if proposal has been uploaded
	if pengajuan.FileProposal == "" {
		return nil, errors.New("proposal belum diupload")
	}

	// 4. Check if status allows assignment
	if pengajuan.StatusProposal != "PENDING" && pengajuan.StatusProposal != "ON_REVIEW" {
		return nil, errors.New("reviewer hanya dapat di-assign untuk proposal dengan status PENDING atau ON_REVIEW")
	}

	// 5. Update pengajuan
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"id_pegawai_reviewer_proposal": idPegawai,
		"status_proposal":              "ON_REVIEW",
		"user_update":                  userUpdateStr,
	}

	if err := database.DB.Model(&pengajuan).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 6. Create plotting record
	plotting := &models.PlottingReviewer{
		IDPengajuan: pengajuan.ID,
		IDPegawai:   idPegawai,
		Tipe:        "PROPOSAL",
	}

	database.DB.Create(plotting)

	// 7. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// ========================================
// ADMIN - ANNOUNCE FINAL RESULT
// ========================================

// AnnounceFinalResult announces final result (LOLOS/TIDAK_LOLOS)
func (s *PengajuanService) AnnounceFinalResult(idPengajuan int, statusFinal string, userID int) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Validate both judul and proposal are reviewed (ACC)
	if pengajuan.StatusJudul != "ACC" {
		return nil, errors.New("judul harus ACC sebelum pengumuman final")
	}

	if pengajuan.StatusProposal != "ACC" {
		return nil, errors.New("proposal harus ACC sebelum pengumuman final")
	}

	// 3. Update status final
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"status_final": statusFinal,
		"user_update":  userUpdateStr,
	}

	if err := database.DB.Model(&pengajuan).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 4. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// ========================================
// REVIEWER - GET MY ASSIGNMENTS
// ========================================

// GetMyAssignments gets all pengajuan assigned to reviewer (pegawai)
func (s *PengajuanService) GetMyAssignments(userID int, tipeFilter string) ([]response.PengajuanListResponse, error) {
	// Note: userID here is from db_user
	// We need to find corresponding pegawai.id from SIMPEG

	// For now, assume userID directly maps to pegawai.id
	// In production, you might need to join db_user with pegawai table
	idPegawai := userID

	// Build query based on tipe filter
	var pengajuanList []models.Pengajuan
	query := database.DB.Where("hapus = ?", 0)

	switch tipeFilter {
	case "JUDUL":
		query = query.Where("id_pegawai_reviewer_judul = ?", idPegawai)
	case "PROPOSAL":
		query = query.Where("id_pegawai_reviewer_proposal = ?", idPegawai)
	default: // "all"
		query = query.Where("id_pegawai_reviewer_judul = ? OR id_pegawai_reviewer_proposal = ?", idPegawai, idPegawai)
	}

	if err := query.Order("tgl_insert DESC").Find(&pengajuanList).Error; err != nil {
		return nil, err
	}

	// Build response list
	result := make([]response.PengajuanListResponse, 0)
	for _, pengajuan := range pengajuanList {
		// Get kategori
		var kategori models.KategoriPKM
		database.DB.Where("id = ?", pengajuan.IDKategori).First(&kategori)

		// Get mahasiswa ketua
		ketua, _ := s.externalService.GetMahasiswaByNIM(pengajuan.NIMKetua)

		// Count anggota
		var anggotaCount int64
		database.DB.Model(&models.PengajuanAnggota{}).
			Where("id_pengajuan = ? AND hapus = ?", pengajuan.ID, 0).
			Count(&anggotaCount)

		// Get reviewer proposal (if assigned)
		var reviewerProposal *external.Pegawai
		if pengajuan.IDPegawaiReviewerProposal != nil {
			reviewerProposal, _ = s.externalService.GetPegawaiByID(*pengajuan.IDPegawaiReviewerProposal)
		}

		// Map to list response
		listResp := s.mapper.MapPengajuanToListResponse(
			&pengajuan,
			ketua,
			&kategori,
			int(anggotaCount),
			reviewerProposal,
		)

		result = append(result, *listResp)
	}

	return result, nil
}

// ========================================
// REVIEWER - REVIEW JUDUL
// ========================================

// ReviewJudul submits review for PKM title
func (s *PengajuanService) ReviewJudul(idPengajuan int, req *request.ReviewJudulRequest, userID int) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Verify reviewer is assigned
	idPegawai := userID // Assume userID maps to pegawai.id
	if pengajuan.IDPegawaiReviewerJudul == nil || *pengajuan.IDPegawaiReviewerJudul != idPegawai {
		return nil, errors.New("anda tidak memiliki akses untuk mereview pengajuan ini")
	}

	// 3. Check if status allows review (must be ON_REVIEW)
	if pengajuan.StatusJudul != "ON_REVIEW" {
		return nil, errors.New("pengajuan harus dalam status ON_REVIEW untuk dapat direview")
	}

	// 4. Get status review info
	var statusReview models.StatusReview
	if err := database.DB.Where("id = ?", req.IDStatusReview).First(&statusReview).Error; err != nil {
		return nil, errors.New("status review tidak valid")
	}

	// 5. START TRANSACTION
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 6. Create review record
	now := time.Now()
	review := &models.ReviewJudul{
		IDPengajuan:    pengajuan.ID,
		IDPegawai:      idPegawai,
		IDStatusReview: req.IDStatusReview,
		Catatan:        req.Catatan,
		TglReview:      &now,
	}

	if err := tx.Create(review).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 7. Update pengajuan status
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"status_judul":         statusReview.KodeStatus, // ACC, REVISI, or TOLAK
		"catatan_review_judul": req.Catatan,
		"tgl_review_judul":     &now,
		"user_update":          userUpdateStr,
	}

	if err := tx.Model(&pengajuan).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 8. Update plotting status to REVIEWED
	tx.Model(&models.PlottingReviewer{}).
		Where("id_pengajuan = ? AND id_pegawai = ? AND tipe = ?", pengajuan.ID, idPegawai, "JUDUL").
		Update("status", 2) // 2 = REVIEWED

	// 9. COMMIT
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 10. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// ========================================
// REVIEWER - REVIEW PROPOSAL
// ========================================

// ReviewProposal submits review for PKM proposal
func (s *PengajuanService) ReviewProposal(idPengajuan int, req *request.ReviewProposalRequest, userID int) (*response.PengajuanResponse, error) {
	// 1. Get pengajuan
	var pengajuan models.Pengajuan
	if err := database.DB.Where("id = ? AND hapus = ?", idPengajuan, 0).First(&pengajuan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}

	// 2. Verify reviewer is assigned
	idPegawai := userID // Assume userID maps to pegawai.id
	if pengajuan.IDPegawaiReviewerProposal == nil || *pengajuan.IDPegawaiReviewerProposal != idPegawai {
		return nil, errors.New("anda tidak memiliki akses untuk mereview pengajuan ini")
	}

	// 3. Check if status allows review (must be ON_REVIEW)
	if pengajuan.StatusProposal != "ON_REVIEW" {
		return nil, errors.New("proposal harus dalam status ON_REVIEW untuk dapat direview")
	}

	// 4. Get status review info
	var statusReview models.StatusReview
	if err := database.DB.Where("id = ?", req.IDStatusReview).First(&statusReview).Error; err != nil {
		return nil, errors.New("status review tidak valid")
	}

	// 5. START TRANSACTION
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 6. Create review record
	now := time.Now()
	review := &models.ReviewProposal{
		IDPengajuan:    pengajuan.ID,
		IDPegawai:      idPegawai,
		IDStatusReview: req.IDStatusReview,
		Catatan:        req.Catatan,
		TglReview:      &now,
	}

	if err := tx.Create(review).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 7. Update pengajuan status
	userUpdateStr := fmt.Sprintf("%d", userID)

	updates := map[string]interface{}{
		"status_proposal":         statusReview.KodeStatus, // ACC, REVISI, or TOLAK
		"catatan_review_proposal": req.Catatan,
		"tgl_review_proposal":     &now,
		"user_update":             userUpdateStr,
	}

	if err := tx.Model(&pengajuan).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 8. Update plotting status to REVIEWED
	tx.Model(&models.PlottingReviewer{}).
		Where("id_pengajuan = ? AND id_pegawai = ? AND tipe = ?", pengajuan.ID, idPegawai, "PROPOSAL").
		Update("status", 2) // 2 = REVIEWED

	// 9. COMMIT
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 10. Return updated detail
	return s.GetPengajuanDetail(idPengajuan)
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// convertToAnggotaModels converts request anggota to models for validation
func (s *PengajuanService) convertToAnggotaModels(anggota []request.AnggotaRequest) []models.PengajuanAnggota {
	result := make([]models.PengajuanAnggota, len(anggota))
	for i, a := range anggota {
		result[i] = models.PengajuanAnggota{
			NIMAnggota: a.NIM,
			IsKetua:    a.IsKetua,
			Urutan:     a.Urutan,
		}
	}
	return result
}
