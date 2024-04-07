package api

import (
	"banner_service/internal/entity"
	"banner_service/internal/repository/postgres"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getUserBanner(ctx *gin.Context) {
	//userBanner := &entity.UserBannerRequest{}
	tagID, err := strconv.ParseUint(ctx.Query("tag_id"), 10, 64)
	if err != nil {
		h.logger.Error("can't get tag id", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	featureID, err := strconv.ParseUint(ctx.Query("feature_id"), 10, 64)
	if err != nil {
		h.logger.Error("can't get feature id", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	useLastVerson, err := strconv.ParseBool(ctx.Query("use_last_revision"))
	if err != nil {
		h.logger.Error("can't get use last version", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userBannerResponse, errGetUserBanner := h.usecase.GetUserBanner(ctx, tagID, featureID, useLastVerson)
	if errGetUserBanner != nil {
		if errors.Is(errGetUserBanner, postgres.ErrorUserBannerNotFound) {
			writeErrorResponse(ctx, http.StatusNotFound, errGetUserBanner.Error())
			return
		}
		h.logger.Error("ошибка получения баннера:", errGetUserBanner)
		writeErrorResponse(ctx, http.StatusInternalServerError, errGetUserBanner.Error())
		return
	}

	h.logger.Info("user banner received", userBannerResponse)
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
		writeErrorResponse(ctx, http.StatusBadRequest, errDecode.Error())
		return
	}

	bannerID, errBannerCreate := h.usecase.CreateBanner(ctx, newBanner)
	if errBannerCreate != nil {
		h.logger.Error("can't create banner", errBannerCreate)
		writeErrorResponse(ctx, http.StatusInternalServerError, errBannerCreate.Error())
		return
	}

	h.logger.Info("banner created", bannerID)
	writePositiveResponse(ctx, http.StatusCreated, bannerID)
	return
}

func (h *Handler) updateBanner(ctx *gin.Context) {
	return
}

func (h *Handler) deleteBanner(ctx *gin.Context) {
	return
}
