package dto

// PaginationRequest untuk query params pagination
type PaginationRequest struct {
	Page  int `json:"page"`  // Default: 1
	Limit int `json:"limit"` // Default: 10, Max: 100
}

// PaginationMeta metadata untuk pagination response
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse response dengan pagination
type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

