package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"rires-be/config"
	"rires-be/internal/dto/response"
)

// APIService handles communication with external API
type APIService struct {
	client *http.Client
}

// NewAPIService creates a new API service instance
func NewAPIService() *APIService {
	return &APIService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// MahasiswaLogin calls external API for mahasiswa login
func (s *APIService) MahasiswaLogin(username, password string) (*response.MahasiswaLoginResponse, error) {
	// Build URL: {baseURL}/mahasiswa/login/{token}/{username}/{password}
	url := fmt.Sprintf("%s/mahasiswa/login/%s/%s/%s",
		config.AppConfig.APIBaseURL,
		config.AppConfig.APIToken,
		username,
		password,
	)

	// Make HTTP request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call mahasiswa API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", string(body))
	}

	// Parse response - handle status as number (1/0) or boolean
	var result struct {
		Status  interface{}                        `json:"status"` // Can be bool or number
		Kode    string                             `json:"kode"`
		Message string                             `json:"message"`
		Data    []response.MahasiswaLoginResponse  `json:"data"` // Array of data
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check status - handle both boolean and number
	statusOK := false
	switch v := result.Status.(type) {
	case bool:
		statusOK = v
	case float64:
		statusOK = v == 1
	case int:
		statusOK = v == 1
	case string:
		statusOK = v == "1" || v == "true"
	}

	if !statusOK {
		return nil, fmt.Errorf("login failed: %s", result.Message)
	}

	// Check if data array is empty
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("login failed: user not found")
	}

	// Return first element from data array
	return &result.Data[0], nil
}

// PegawaiLogin calls external API for pegawai login
func (s *APIService) PegawaiLogin(username, password string) (*response.PegawaiLoginResponse, error) {
	// Build URL: {baseURL}/pegawai/login/{token}/{username}/{password}
	url := fmt.Sprintf("%s/pegawai/login/%s/%s/%s",
		config.AppConfig.APIBaseURL,
		config.AppConfig.APIToken,
		username,
		password,
	)

	// Make HTTP request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call pegawai API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", string(body))
	}

	// Parse response - handle status as number (1/0) or boolean
	var result struct {
		Status  interface{}                       `json:"status"` // Can be bool or number
		Kode    string                            `json:"kode"`
		Message string                            `json:"message"`
		Data    []response.PegawaiLoginResponse   `json:"data"` // Array of data
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check status - handle both boolean and number
	statusOK := false
	switch v := result.Status.(type) {
	case bool:
		statusOK = v
	case float64:
		statusOK = v == 1
	case int:
		statusOK = v == 1
	case string:
		statusOK = v == "1" || v == "true"
	}

	if !statusOK {
		return nil, fmt.Errorf("login failed: %s", result.Message)
	}

	// Check if data array is empty
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("login failed: user not found")
	}

	// Return first element from data array
	return &result.Data[0], nil
}