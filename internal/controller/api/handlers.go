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
	newBanner := &entity.Banner{}
	errDecode := json.NewDecoder(ctx.Request.Body).Decode(&newBanner)
	if errDecode != nil {
		h.logger.Error("can't decode banner", errDecode)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errDecode.Error()})
		return
	}

	bannerID, errBannerCreate := h.usecase.CreateBanner(ctx, newBanner)
	if errBannerCreate != nil {
		h.logger.Error("can't create banner", errBannerCreate)
		ctx.JSON(http.StatusInternalServerError, errBannerCreate)
		return
	}

	h.logger.Info("banner created", bannerID)
	ctx.JSON(http.StatusCreated, bannerID)
	return
}

func (h *Handler) updateBanner(ctx *gin.Context) {
	return
}

func (h *Handler) deleteBanner(ctx *gin.Context) {
	return
}
