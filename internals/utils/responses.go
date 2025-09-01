package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func APISuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

func APIError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
	c.Abort()
}

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func NewPagination(page, limit int, total int64) *Pagination {
	totalPage := int(total) / limit
	if int(total)%limit > 0 {
		totalPage++
	}

	return &Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}
}