package services

import (
	"encoding/json"
	"rires-be/internal/dto/response"
	"rires-be/internal/models"
	"rires-be/internal/models/external"
)

// MapperService handles conversion between models and DTOs
type MapperService struct{}

// NewMapperService creates a new mapper service
func NewMapperService() *MapperService {
	return &MapperService{}
}

// ========================================
// MAHASISWA MAPPING
// ========================================

// MapMahasiswaToResponse converts external.Mahasiswa to response.MahasiswaResponse
func (m *MapperService) MapMahasiswaToResponse(mhs *external.Mahasiswa, anggota *models.PengajuanAnggota) *response.MahasiswaResponse {
	if mhs == nil {
		return nil
	}

	resp := &response.MahasiswaResponse{
		NIM:      mhs.KodeSiswa,
		Nama:     mhs.NamaSiswa,
		HP:       mhs.HPSiswa,
		Angkatan: mhs.TahunMasuk,
	}

	// Add team context if anggota data provided
	if anggota != nil {
		resp.IsKetua = &anggota.IsKetua
	}

	// Map Prodi if available
	if mhs.Prodi != nil {
		resp.Prodi = m.MapProdiToResponse(mhs.Prodi)
	}

	return resp
}

// MapProdiToResponse converts external.Prodi to response.ProdiResponse
func (m *MapperService) MapProdiToResponse(prodi *external.Prodi) *response.ProdiResponse {
	if prodi == nil {
		return nil
	}

	resp := &response.ProdiResponse{
		Kode:        prodi.Kode,
		NamaProdi:   prodi.NamaDepart,
		NamaSingkat: prodi.NamaSingkat,
	}

	// Map Fakultas if available
	if prodi.Fakultas != nil {
		resp.Fakultas = m.MapFakultasToResponse(prodi.Fakultas)
	}

	return resp
}

// MapFakultasToResponse converts external.Fakultas to response.FakultasResponse
func (m *MapperService) MapFakultasToResponse(fakultas *external.Fakultas) *response.FakultasResponse {
	if fakultas == nil {
		return nil
	}

	return &response.FakultasResponse{
		Kode:         fakultas.Kode,
		NamaFakultas: fakultas.NamaFakultas,
		NamaSingkat:  fakultas.NamaFakPendek,
	}
}

// ========================================
// PEGAWAI MAPPING
// ========================================

// MapPegawaiToResponse converts external.Pegawai to response.PegawaiResponse
func (m *MapperService) MapPegawaiToResponse(pegawai *external.Pegawai) *response.PegawaiResponse {
	if pegawai == nil {
		return nil
	}

	resp := &response.PegawaiResponse{
		ID:          pegawai.ID,
		Nama:        pegawai.NamaPegawai,
		NamaLengkap: pegawai.GetNamaLengkap(),
		Email:       pegawai.Email,
		HP:          pegawai.HP,
	}

	return resp
}

// ========================================
// KATEGORI MAPPING
// ========================================

// MapKategoriToResponse converts models.KategoriPKM to response.KategoriResponse
func (m *MapperService) MapKategoriToResponse(kategori *models.KategoriPKM) *response.KategoriResponse {
	if kategori == nil {
		return nil
	}

	return &response.KategoriResponse{
		ID:           kategori.ID,
		NamaKategori: kategori.NamaKategori,
	}
}

// ========================================
// PARAMETER MAPPING
// ========================================

// MapParameterToResponse converts models.ParameterPKM to response.ParameterResponse
func (m *MapperService) MapParameterToResponse(param *models.ParameterPKM) *response.ParameterResponse {
	if param == nil {
		return nil
	}

	resp := &response.ParameterResponse{
		ID:          param.ID,
		IDParameter: param.IDParameter,
		Nilai:       param.Nilai,
	}

	// Add label and tipe_input from ParameterForm if available
	if param.ParameterForm != nil {
		resp.Label = param.ParameterForm.Label
		resp.TipeInput = param.ParameterForm.TipeInput
	}

	return resp
}

// ========================================
// REVIEW MAPPING
// ========================================

// MapReviewJudulToResponse converts models.ReviewJudul to response.ReviewResponse
func (m *MapperService) MapReviewJudulToResponse(review *models.ReviewJudul, pegawai *external.Pegawai) *response.ReviewResponse {
	if review == nil {
		return nil
	}

	resp := &response.ReviewResponse{
		ID:         review.ID,
		TipeReview: "JUDUL",
		Catatan:    review.Catatan,
		TglReview:  review.TglReview,
	}

	// Map status review
	if review.StatusReview != nil {
		resp.StatusReview = review.StatusReview.KodeStatus
	}

	// Map reviewer
	if pegawai != nil {
		resp.Reviewer = m.MapPegawaiToResponse(pegawai)
	}

	return resp
}

// MapReviewProposalToResponse converts models.ReviewProposal to response.ReviewResponse
func (m *MapperService) MapReviewProposalToResponse(review *models.ReviewProposal, pegawai *external.Pegawai) *response.ReviewResponse {
	if review == nil {
		return nil
	}

	resp := &response.ReviewResponse{
		ID:         review.ID,
		TipeReview: "PROPOSAL",
		Catatan:    review.Catatan,
		TglReview:  review.TglReview,
	}

	// Map status review
	if review.StatusReview != nil {
		resp.StatusReview = review.StatusReview.KodeStatus
	}

	// Map reviewer
	if pegawai != nil {
		resp.Reviewer = m.MapPegawaiToResponse(pegawai)
	}

	return resp
}

// ========================================
// PENGAJUAN MAPPING
// ========================================

// MapPengajuanToListResponse converts models.Pengajuan to response.PengajuanListResponse
func (m *MapperService) MapPengajuanToListResponse(
	pengajuan *models.Pengajuan,
	ketua *external.Mahasiswa,
	kategori *models.KategoriPKM,
	jumlahAnggota int,
	reviewerProposal *external.Pegawai,
	reviewerJudulNama string,
) *response.PengajuanListResponse {
	if pengajuan == nil {
		return nil
	}

	resp := &response.PengajuanListResponse{
		ID:             pengajuan.ID,
		KodePengajuan:  pengajuan.KodePengajuan,
		Judul:          pengajuan.Judul,
		Tahun:          pengajuan.Tahun,
		StatusJudul:    pengajuan.StatusJudul,
		StatusProposal: pengajuan.StatusProposal,
		StatusFinal:    pengajuan.StatusFinal,
		JumlahAnggota:  jumlahAnggota,
		TglInsert:      pengajuan.TglInsert,

		// Flat fields
		NIMKetua:        pengajuan.NIMKetua,
		EmailKetua:      pengajuan.EmailKetua,
		NoHPKetua:       pengajuan.NoHPKetua,
		ProgramStudi:    pengajuan.ProgramStudi,
		Fakultas:        pengajuan.Fakultas,
		CatatanProposal: pengajuan.CatatanReviewProposal,
		TanggalReview:   pengajuan.TglReviewProposal,
		FileProposal:    pengajuan.FileProposal,
	}

	// Map kategori
	if kategori != nil {
		resp.Kategori = m.MapKategoriToResponse(kategori)
		resp.NamaKategori = kategori.NamaKategori
	}

	// Map ketua
	if ketua != nil {
		resp.Ketua = m.MapMahasiswaToResponse(ketua, nil)
		resp.NamaKetua = ketua.NamaSiswa
	}

	// Set nama_reviewer - prioritize judul reviewer, fallback to proposal reviewer
	if reviewerJudulNama != "" {
		resp.NamaReviewer = reviewerJudulNama
	} else if reviewerProposal != nil {
		resp.NamaReviewer = reviewerProposal.GetNamaLengkap()
	}

	return resp
}

// MapPengajuanToDetailResponse converts models.Pengajuan to response.PengajuanResponse (full detail)
func (m *MapperService) MapPengajuanToDetailResponse(
	pengajuan *models.Pengajuan,
	ketua *external.Mahasiswa,
	anggotaList []external.Mahasiswa,
	anggotaModels []models.PengajuanAnggota,
	kategori *models.KategoriPKM,
	parameterList []models.ParameterPKM,
	reviewerJudul *external.Pegawai,
	reviewerProposal *external.Pegawai,
	reviewJudulHistory []models.ReviewJudul,
	reviewProposalHistory []models.ReviewProposal,
) *response.PengajuanResponse {
	if pengajuan == nil {
		return nil
	}

	resp := &response.PengajuanResponse{
		ID:            pengajuan.ID,
		KodePengajuan: pengajuan.KodePengajuan,
		Judul:         pengajuan.Judul,
		Tahun:         pengajuan.Tahun,
		// Biodata Ketua
		NamaKetua:       pengajuan.NamaKetua,
		NIMKetua:        pengajuan.NIMKetua,
		EmailKetua:      pengajuan.EmailKetua,
		NoHPKetua:       pengajuan.NoHPKetua,
		ProgramStudi:    pengajuan.ProgramStudi,
		Fakultas:        pengajuan.Fakultas,
		DosenPembimbing: pengajuan.DosenPembimbing,
		TglPengajuan:    pengajuan.TglPengajuan,
		// Kategori
		IDKategori: pengajuan.IDKategori,
		// Status
		StatusJudul:           pengajuan.StatusJudul,
		StatusProposal:        pengajuan.StatusProposal,
		StatusFinal:           pengajuan.StatusFinal,
		FileProposal:          pengajuan.FileProposal,
		CatatanReviewJudul:    pengajuan.CatatanReviewJudul,
		TglReviewJudul:        pengajuan.TglReviewJudul,
		CatatanReviewProposal: pengajuan.CatatanReviewProposal,
		TglReviewProposal:     pengajuan.TglReviewProposal,
		TglInsert:             pengajuan.TglInsert,
		TglUpdate:             pengajuan.TglUpdate,
	}

	// Map anggota_list from local DB
	if len(anggotaModels) > 0 {
		resp.AnggotaList = make([]response.AnggotaResponse, len(anggotaModels))
		for i, anggota := range anggotaModels {
			resp.AnggotaList[i] = response.AnggotaResponse{
				ID:          anggota.ID,
				NIMAnggota:  anggota.NIMAnggota,
				NamaAnggota: anggota.NamaAnggota,
				IsKetua:     anggota.IsKetua,
				Urutan:      anggota.Urutan,
			}
		}
	}

	// Map kategori
	if kategori != nil {
		resp.Kategori = m.MapKategoriToResponse(kategori)
		resp.NamaKategori = kategori.NamaKategori // Flat field
	}

	// Map ketua
	if ketua != nil {
		// Find ketua in anggotaModels to get IsKetua flag
		var ketuaModel *models.PengajuanAnggota
		for _, anggota := range anggotaModels {
			if anggota.NIMAnggota == ketua.KodeSiswa {
				ketuaModel = &anggota
				break
			}
		}
		resp.Ketua = m.MapMahasiswaToResponse(ketua, ketuaModel)
	}

	// Map anggota
	resp.Anggota = make([]response.MahasiswaResponse, 0)
	for _, mhs := range anggotaList {
		// Find corresponding anggota model
		var anggotaModel *models.PengajuanAnggota
		for _, anggota := range anggotaModels {
			if anggota.NIMAnggota == mhs.KodeSiswa {
				anggotaModel = &anggota
				break
			}
		}

		if anggotaResp := m.MapMahasiswaToResponse(&mhs, anggotaModel); anggotaResp != nil {
			resp.Anggota = append(resp.Anggota, *anggotaResp)
		}
	}

	// Parse parameter_data JSON string to map
	if pengajuan.ParameterData != "" {
		var paramData map[string]interface{}
		if err := json.Unmarshal([]byte(pengajuan.ParameterData), &paramData); err == nil {
			resp.ParameterData = paramData
		}
	}

	// Map reviewers
	if reviewerJudul != nil {
		resp.ReviewerJudul = m.MapPegawaiToResponse(reviewerJudul)
	}
	if reviewerProposal != nil {
		resp.ReviewerProposal = m.MapPegawaiToResponse(reviewerProposal)
	}

	// Map review history judul
	resp.ReviewJudulHistory = make([]response.ReviewResponse, 0)
	for _, review := range reviewJudulHistory {
		if reviewResp := m.MapReviewJudulToResponse(&review, reviewerJudul); reviewResp != nil {
			resp.ReviewJudulHistory = append(resp.ReviewJudulHistory, *reviewResp)
		}
	}

	// Map review history proposal
	resp.ReviewProposalHistory = make([]response.ReviewResponse, 0)
	for _, review := range reviewProposalHistory {
		if reviewResp := m.MapReviewProposalToResponse(&review, reviewerProposal); reviewResp != nil {
			resp.ReviewProposalHistory = append(resp.ReviewProposalHistory, *reviewResp)
		}
	}

	return resp
}
