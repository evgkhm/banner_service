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
		if errors.Is(errGetUserBanner, postgres.ErrUserBanner) {
			writeErrorResponse(ctx, http.StatusNotFound, errGetUserBanner.Error())
			return
		}
		h.logger.Error("ошибка получения баннера:", errGetUserBanner)
		writeErrorResponse(ctx, http.StatusInternalServerError, errGetUserBanner.Error())
		return
	}

	h.logger.Info("user banner received", userBannerResponse)
	ctx.JSON(http.StatusOK, userBannerResponse)
}

func (h *Handler) getBanners(ctx *gin.Context) {
	tagIDs, errParseTagIDs := strconv.ParseUint(ctx.Query("tag_id"), 10, 64)
	if errParseTagIDs != nil {
		h.logger.Error("can't parse tag ids", errParseTagIDs)
		writeErrorResponse(ctx, http.StatusBadRequest, errParseTagIDs.Error())
		return
	}

	featureID, err := strconv.ParseUint(ctx.Query("feature_id"), 10, 64)
	if err != nil {
		h.logger.Error("can't get feature id", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	limit, err := strconv.ParseUint(ctx.Query("limit"), 10, 64)
	if err != nil {
		h.logger.Error("can't get limit", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.ParseUint(ctx.Query("offset"), 10, 64)
	if err != nil {
		h.logger.Error("can't get offset", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	usersBanners, errGetBanners := h.usecase.GetBanners(ctx, tagIDs, featureID, limit, offset)
	if errGetBanners != nil {
		h.logger.Error("ошибка получения баннеров:", errGetBanners)
		writeErrorResponse(ctx, http.StatusInternalServerError, errGetBanners.Error())
		return
	}
	h.logger.Info("полученные баннеры", usersBanners)
	ctx.JSON(http.StatusOK, usersBanners)
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
}

func (h *Handler) updateBanner(ctx *gin.Context) {
	bannerID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("can't get banner id", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	newBanner := &entity.Banner{}
	errDecode := json.NewDecoder(ctx.Request.Body).Decode(&newBanner)
	if errDecode != nil {
		h.logger.Error("can't decode banner", errDecode)
		writeErrorResponse(ctx, http.StatusBadRequest, errDecode.Error())
		return
	}

	errUpdate := h.usecase.UpdateBanner(ctx, bannerID, newBanner)
	if errUpdate != nil {
		h.logger.Error("can't update banner", errUpdate)
		writeErrorResponse(ctx, http.StatusInternalServerError, errUpdate.Error())
		return
	}
	h.logger.Info("banner updated", bannerID)
	ctx.JSON(http.StatusOK, "OK")
}

func (h *Handler) deleteBanner(ctx *gin.Context) {
	bannerID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("can't get banner id", err)
		writeErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	errDelete := h.usecase.DeleteBanner(ctx, bannerID)
	if errDelete != nil {
		h.logger.Error("can't delete banner", errDelete)
		writeErrorResponse(ctx, http.StatusInternalServerError, errDelete.Error())
		return
	}
	h.logger.Info("Баннер успешно удален")
	ctx.JSON(http.StatusNoContent, "Баннер успешно удален")
}
