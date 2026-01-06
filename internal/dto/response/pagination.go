package response

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

// PaginatedResponse represents paginated data response
type PaginatedResponse struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

// NewPaginationResponse creates pagination metadata
func NewPaginationResponse(page, perPage int, totalRecords int64) *PaginationResponse {
	totalPages := int((totalRecords + int64(perPage) - 1) / int64(perPage))
	
	return &PaginationResponse{
		CurrentPage:  page,
		PerPage:      perPage,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}