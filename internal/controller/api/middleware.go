package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func middlewareForAdmin() gin.HandlerFunc {
	requiredAdminToken := os.Getenv("API_TOKEN")

	if requiredAdminToken == "" {
		log.Fatal("Please set API_TOKEN environment variable")
	}

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			respondWithError(c, 401, "Token required")
			return
		}

		if token != requiredAdminToken {
			respondWithError(c, 401, "Invalid token")
			return
		}
	}
}

func middlewareForAdminOrUser() gin.HandlerFunc {
	requiredUserToken := os.Getenv("API_USER_TOKEN")
	requiredAdminToken := os.Getenv("API_TOKEN")

	if requiredUserToken == "" || requiredAdminToken == "" {
		log.Fatal("Please set API_TOKEN environment variable")
	}

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			respondWithError(c, 401, "Token required")
			return
		}

		if token != requiredUserToken && token != requiredAdminToken {
			respondWithError(c, 401, "Invalid token")
			return
		}
	}
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
