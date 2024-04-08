package api

import (
	"banner_service/internal/entity"
	"banner_service/internal/repository/postgres"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) getUserBanner(ctx *gin.Context) {
	//authMiddleware := DummyMiddleware()
	//authMiddleware(ctx)

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

func (h *Handler) getBanners(ctx *gin.Context) {
	//tagIDs, errParseTagIDs := parseTagIDs(ctx.Query("tag_ids"))
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
	return
}

func parseTagIDs(tagIDsParam string) ([]uint64, error) {
	var tagIDs []uint64

	ids := strings.Split(tagIDsParam, ",")
	for _, id := range ids {
		tagID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}

	return tagIDs, nil
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
	return
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
	return
}
