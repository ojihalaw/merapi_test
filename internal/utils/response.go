package utils

import "github.com/gofiber/fiber/v2"

type PaginationRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	OrderBy string `json:"order_by"`
	SortBy  string `json:"sort_by"`
	Search  string `json:"search"`
}

type PaginationResponse struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	OrderBy   string `json:"order_by"`
	SortBy    string `json:"sort_by"`
	Search    string `json:"search"`
	TotalData int64  `json:"total_data"`
	TotalPage int    `json:"total_page"`
}

func DefaultSuccessResponse(code int, message string) fiber.Map {
	return fiber.Map{
		"code":    code,
		"status":  true,
		"message": message,
	}
}

func SuccessResponse(code int, message string, data interface{}) fiber.Map {
	return fiber.Map{
		"code":    code,
		"status":  true,
		"message": message,
		"result":  data,
	}
}

func SuccessResponseWithPagination(code int, message string, data interface{}, pagination *PaginationResponse) fiber.Map {
	return fiber.Map{
		"code":       code,
		"status":     true,
		"message":    message,
		"result":     data,
		"pagination": pagination,
	}
}

func ErrorResponse(code int, message string) fiber.Map {
	return fiber.Map{
		"code":    code,
		"status":  false,
		"message": message,
	}
}
