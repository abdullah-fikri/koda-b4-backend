package models

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Pagination any `json:"pagination,omitempty"`
	Data    any    `json:"data,omitempty"`
}
