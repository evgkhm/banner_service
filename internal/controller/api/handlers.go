package api

import (
	"banner_service/internal/entity"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getUserBanner(ctx *gin.Context) {
	userBanner := &entity.UserBannerRequest{}
	errDecode := json.NewDecoder(ctx.Request.Body).Decode(userBanner)
	if errDecode != nil {
		h.logger.Error("can't decode user banner", errDecode)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errDecode.Error()})
		return
	}

	userBannerResponse, errGetUserBanner := h.usecase.GetUserBanner(ctx, userBanner)
	if errGetUserBanner != nil {
		h.logger.Error("can't get user banner", errGetUserBanner)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errGetUserBanner})
		return
	}

	ctx.JSON(http.StatusOK, userBannerResponse)
	return
}

func (h *Handler) getBanner(ctx *gin.Context) {
	return
}

func (h *Handler) createBanner(ctx *gin.Context) {
	return
}

func (h *Handler) updateBanner(ctx *gin.Context) {
	return
}

func (h *Handler) deleteBanner(ctx *gin.Context) {
	return
}
