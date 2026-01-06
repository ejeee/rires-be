package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileUploadService handles file upload operations
type FileUploadService struct {
	UploadDir string
	MaxSize   int64 // in bytes
}

// NewFileUploadService creates a new instance of FileUploadService
func NewFileUploadService() *FileUploadService {
	return &FileUploadService{
		UploadDir: "./uploads/proposals", // Default upload directory
		MaxSize:   2.5 * 1024 * 1024,     // 2.5 MB in bytes
	}
}

// UploadProposal uploads proposal file and returns the file path
func (s *FileUploadService) UploadProposal(file *multipart.FileHeader, kodePengajuan string) (string, error) {
	// Validate file
	if err := s.ValidateFile(file); err != nil {
		return "", err
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(s.UploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	filename := s.GenerateFilename(kodePengajuan, file.Filename)
	filepath := filepath.Join(s.UploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path
	return filename, nil
}

// ValidateFile validates uploaded file (extension, size, mime type)
func (s *FileUploadService) ValidateFile(file *multipart.FileHeader) error {
	// Check if file is nil
	if file == nil {
		return errors.New("no file uploaded")
	}

	// Validate file size
	if file.Size > s.MaxSize {
		return fmt.Errorf("file size exceeds maximum limit of %.1f MB", float64(s.MaxSize)/(1024*1024))
	}

	if file.Size == 0 {
		return errors.New("uploaded file is empty")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := []string{".pdf", ".doc", ".docx"}
	
	isValid := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid file extension: %s. Allowed: PDF, DOC, DOCX", ext)
	}

	return nil
}

// GenerateFilename generates unique filename for proposal
// Format: proposal_{kodePengajuan}_{timestamp}.{ext}
func (s *FileUploadService) GenerateFilename(kodePengajuan, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("proposal_%s_%s%s", kodePengajuan, timestamp, ext)
}

// DeleteFile deletes a file from the upload directory
func (s *FileUploadService) DeleteFile(filename string) error {
	if filename == "" {
		return nil // Nothing to delete
	}

	filepath := filepath.Join(s.UploadDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	// Delete file
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFilePath returns full path to uploaded file
func (s *FileUploadService) GetFilePath(filename string) string {
	return filepath.Join(s.UploadDir, filename)
}

// FileExists checks if file exists in upload directory
func (s *FileUploadService) FileExists(filename string) bool {
	if filename == "" {
		return false
	}

	filepath := filepath.Join(s.UploadDir, filename)
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// GetFileSize returns file size in bytes
func (s *FileUploadService) GetFileSize(filename string) (int64, error) {
	filepath := filepath.Join(s.UploadDir, filename)
	
	info, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}