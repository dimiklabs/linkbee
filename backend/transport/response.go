package transport

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/shafikshaon/shortlink/logger"
)

type StandardResponse struct {
	Timestamp   int64  `json:"timestamp"`
	RequestID   string `json:"request_id"`
	Description string `json:"description,omitempty"`
	Data        any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Timestamp   int64  `json:"timestamp"`
	RequestID   string `json:"request_id"`
	ErrorCode   string `json:"error_code"`
	Description string `json:"description"`
	Data        any    `json:"data,omitempty"`
}

func RespondWithSuccess(c *gin.Context, statusCode int, description string, data any) {
	c.JSON(statusCode, StandardResponse{
		Timestamp:   time.Now().Unix(),
		RequestID:   getRequestID(c),
		Description: description,
		Data:        data,
	})
}

func RespondWithError(c *gin.Context, statusCode int, errorCode, description string) {
	c.JSON(statusCode, ErrorResponse{
		Timestamp:   time.Now().Unix(),
		RequestID:   getRequestID(c),
		ErrorCode:   errorCode,
		Description: description,
	})
}

func RespondWithErrorData(c *gin.Context, statusCode int, errorCode, description string, data any) {
	c.JSON(statusCode, ErrorResponse{
		Timestamp:   time.Now().Unix(),
		RequestID:   getRequestID(c),
		ErrorCode:   errorCode,
		Description: description,
		Data:        data,
	})
}

func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(string(logger.RequestIDKey)); exists {
		return requestID.(string)
	}
	return ""
}
