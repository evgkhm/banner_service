package api

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error" example:"error message"`
}

func writeErrorResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Error: msg})
}
