package services

import (
	"errors"
	"rires-be/internal/models/external"
	"rires-be/pkg/database"
)

// ExternalDataService handles queries to external databases (NEOMAA, NEOMAAREF, SIMPEG)
type ExternalDataService struct{}

// NewExternalDataService creates a new instance of ExternalDataService
func NewExternalDataService() *ExternalDataService {
	return &ExternalDataService{}
}

// ========================================
// MAHASISWA - NEOMAA
// ========================================

// GetMahasiswaByNIM fetches mahasiswa data from NEOMAA by NIM
func (s *ExternalDataService) GetMahasiswaByNIM(nim string) (*external.Mahasiswa, error) {
	if database.DBNeomaa == nil {
		return nil, errors.New("NEOMAA database not connected")
	}

	var mahasiswa external.Mahasiswa
	if err := database.DBNeomaa.Where("nim = ? AND status = ?", nim, 1).First(&mahasiswa).Error; err != nil {
		return nil, err
	}

	return &mahasiswa, nil
}

// GetMahasiswaByNIMs fetches multiple mahasiswa by NIMs
func (s *ExternalDataService) GetMahasiswaByNIMs(nims []string) ([]external.Mahasiswa, error) {
	if database.DBNeomaa == nil {
		return nil, errors.New("NEOMAA database not connected")
	}

	var mahasiswaList []external.Mahasiswa
	if err := database.DBNeomaa.Where("nim IN ? AND status = ?", nims, 1).Find(&mahasiswaList).Error; err != nil {
		return nil, err
	}

	return mahasiswaList, nil
}

// ValidateNIMExists checks if NIM exists and is active in NEOMAA
func (s *ExternalDataService) ValidateNIMExists(nim string) bool {
	if database.DBNeomaa == nil {
		return false
	}

	var count int64
	database.DBNeomaa.Model(&external.Mahasiswa{}).
		Where("nim = ? AND status = ?", nim, 1).
		Count(&count)

	return count > 0
}

// GetMahasiswaWithProdi fetches mahasiswa data with prodi information
func (s *ExternalDataService) GetMahasiswaWithProdi(nim string) (*external.Mahasiswa, error) {
	if database.DBNeomaa == nil {
		return nil, errors.New("NEOMAA database not connected")
	}

	var mahasiswa external.Mahasiswa
	if err := database.DBNeomaa.Preload("Prodi.Fakultas").
		Where("nim = ? AND status = ?", nim, 1).
		First(&mahasiswa).Error; err != nil {
		return nil, err
	}

	return &mahasiswa, nil
}

// ========================================
// PEGAWAI - SIMPEG
// ========================================

// GetPegawaiByID fetches pegawai data from SIMPEG by ID
func (s *ExternalDataService) GetPegawaiByID(id int) (*external.Pegawai, error) {
	if database.DBSimpeg == nil {
		return nil, errors.New("SIMPEG database not connected")
	}

	var pegawai external.Pegawai
	if err := database.DBSimpeg.Where("id = ? AND status = ?", id, 1).First(&pegawai).Error; err != nil {
		return nil, err
	}

	return &pegawai, nil
}

// GetPegawaiByIDs fetches multiple pegawai by IDs
func (s *ExternalDataService) GetPegawaiByIDs(ids []int) ([]external.Pegawai, error) {
	if database.DBSimpeg == nil {
		return nil, errors.New("SIMPEG database not connected")
	}

	var pegawaiList []external.Pegawai
	if err := database.DBSimpeg.Where("id IN ? AND status = ?", ids, 1).Find(&pegawaiList).Error; err != nil {
		return nil, err
	}

	return pegawaiList, nil
}

// ValidatePegawaiExists checks if pegawai ID exists and is active in SIMPEG
func (s *ExternalDataService) ValidatePegawaiExists(id int) bool {
	if database.DBSimpeg == nil {
		return false
	}

	var count int64
	database.DBSimpeg.Model(&external.Pegawai{}).
		Where("id = ? AND status = ?", id, 1).
		Count(&count)

	return count > 0
}

// GetPegawaiWithFakultas fetches pegawai data with fakultas information
func (s *ExternalDataService) GetPegawaiWithFakultas(id int) (*external.Pegawai, error) {
	if database.DBSimpeg == nil {
		return nil, errors.New("SIMPEG database not connected")
	}

	var pegawai external.Pegawai
	if err := database.DBSimpeg.Preload("Fakultas").
		Where("id = ? AND status = ?", id, 1).
		First(&pegawai).Error; err != nil {
		return nil, err
	}

	return &pegawai, nil
}

// GetAllReviewers fetches all pegawai that can be reviewers (active status)
func (s *ExternalDataService) GetAllReviewers() ([]external.Pegawai, error) {
	if database.DBSimpeg == nil {
		return nil, errors.New("SIMPEG database not connected")
	}

	var pegawaiList []external.Pegawai
	if err := database.DBSimpeg.Where("status = ?", 1).
		Order("nama ASC").
		Find(&pegawaiList).Error; err != nil {
		return nil, err
	}

	return pegawaiList, nil
}

// ========================================
// PRODI & FAKULTAS - NEOMAAREF
// ========================================

// GetProdiByID fetches prodi data from NEOMAAREF by ID
func (s *ExternalDataService) GetProdiByID(id int) (*external.Prodi, error) {
	if database.DBNeomaaRef == nil {
		return nil, errors.New("NEOMAAREF database not connected")
	}

	var prodi external.Prodi
	if err := database.DBNeomaaRef.Where("id = ? AND status = ?", id, 1).First(&prodi).Error; err != nil {
		return nil, err
	}

	return &prodi, nil
}

// GetFakultasByID fetches fakultas data from NEOMAAREF by ID
func (s *ExternalDataService) GetFakultasByID(id int) (*external.Fakultas, error) {
	if database.DBNeomaaRef == nil {
		return nil, errors.New("NEOMAAREF database not connected")
	}

	var fakultas external.Fakultas
	if err := database.DBNeomaaRef.Where("id = ? AND status = ?", id, 1).First(&fakultas).Error; err != nil {
		return nil, err
	}

	return &fakultas, nil
}

// GetAllProdi fetches all active prodi
func (s *ExternalDataService) GetAllProdi() ([]external.Prodi, error) {
	if database.DBNeomaaRef == nil {
		return nil, errors.New("NEOMAAREF database not connected")
	}

	var prodiList []external.Prodi
	if err := database.DBNeomaaRef.Preload("Fakultas").
		Where("status = ?", 1).
		Order("nama_prodi ASC").
		Find(&prodiList).Error; err != nil {
		return nil, err
	}

	return prodiList, nil
}

// GetAllFakultas fetches all active fakultas
func (s *ExternalDataService) GetAllFakultas() ([]external.Fakultas, error) {
	if database.DBNeomaaRef == nil {
		return nil, errors.New("NEOMAAREF database not connected")
	}

	var fakultasList []external.Fakultas
	if err := database.DBNeomaaRef.Where("status = ?", 1).
		Order("nama_fakultas ASC").
		Find(&fakultasList).Error; err != nil {
		return nil, err
	}

	return fakultasList, nil
}