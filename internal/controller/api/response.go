package api

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error" example:"error message"`
}

type positiveResponse struct {
	Response uint64 `json:"banner_id" example:"1"`
}

func writeErrorResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Error: msg})
}

func writePositiveResponse(c *gin.Context, statusCode int, response uint64) {
	c.JSON(statusCode, positiveResponse{Response: response})
}
