package response

// APIResponse represents standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse creates a success API response
func SuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error API response
func ErrorResponse(message string, err interface{}) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}

// ValidationErrorResponse represents validation error details
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}