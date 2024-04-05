package api

import "github.com/gin-gonic/gin"

type Handler struct {
	*gin.Engine
	usecase useCase
	logger  logger
}

func New(usecase useCase, log logger) *Handler {
	h := &Handler{
		Engine:  gin.New(),
		usecase: usecase,
		logger:  log,
	}

	h.Use(gin.Recovery())

	// Swagger
	//h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//api := h.Group("/api")
	h.GET("/user_banner", h.userBanner)
	//api.GET("products/get", h.getProducts)

	return h

}
