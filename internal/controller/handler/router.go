package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	*gin.Engine
	usecase usecase
	logger  logger
}

func NewHandler(usecase usecase, log logger) *Handler {
	h := &Handler{
		Engine:  gin.New(),
		usecase: usecase,
		logger:  log,
	}

	h.Use(gin.Recovery())

	// Swagger
	h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := h.Group("/api")

	api.GET("products/get", h.getProducts)

	return h

}
