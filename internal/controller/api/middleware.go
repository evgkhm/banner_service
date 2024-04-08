package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func DummyMiddleware() gin.HandlerFunc {
	requiredAdminToken := os.Getenv("API_TOKEN")
	requiredUserToken := os.Getenv("API_USER_TOKEN")
	// We want to make sure the token is set, bail if not
	if requiredUserToken == "" || requiredAdminToken == "" {
		log.Fatal("Please set API_TOKEN environment variable")
	}

	return func(c *gin.Context) {
		// Проверяем пользовательский токен
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
