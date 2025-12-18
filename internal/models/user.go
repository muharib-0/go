package models

import (
	"time"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UserResponse represents the response body for a user (without age)
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

// UserWithAgeResponse represents the response body for a user (with age)
type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age"`
}

// PaginatedUsersResponse represents a paginated list of users
type PaginatedUsersResponse struct {
	Users      []UserWithAgeResponse `json:"users"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// PaginationQuery represents pagination query parameters
type PaginationQuery struct {
	Page     int `query:"page" validate:"omitempty,min=1"`
	PageSize int `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// CalculateAge calculates the age based on date of birth
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// Check if birthday hasn't occurred this year yet
	if now.YearDay() < dob.YearDay() {
		years--
	}

	// Handle edge case where the calculated age would be negative
	if years < 0 {
		return 0
	}

	return years
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}
