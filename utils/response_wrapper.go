package utils

import "github.com/gin-gonic/gin"

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
	})
}
