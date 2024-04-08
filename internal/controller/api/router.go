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

	h.GET("/user_banner", middlewareForAdminOrUser(), h.getUserBanner)
	h.GET("/banner", middlewareForAdmin(), h.getBanners)
	h.POST("/banner", middlewareForAdmin(), h.createBanner)
	h.PATCH("/banner/:id", middlewareForAdmin(), h.updateBanner)
	h.DELETE("/banner/:id", middlewareForAdmin(), h.deleteBanner)

	return h

}
