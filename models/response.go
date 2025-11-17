package models

import "backend/lib"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Pagination *lib.PaginationData `json:"pagination,omitempty"`
	Data    any    `json:"data,omitempty"`
}
