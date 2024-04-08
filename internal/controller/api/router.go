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

	h.GET("/user_banner", h.getUserBanner)
	h.GET("/banner", h.getBanners)
	h.POST("/banner", h.createBanner)
	h.PATCH("/banner/:id", h.updateBanner)
	h.DELETE("/banner/{id}", h.deleteBanner)

	return h

}
