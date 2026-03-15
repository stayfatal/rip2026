package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

func (h *Handler) AddToSystemLoad(ctx *gin.Context) {
	strategyIDStr := ctx.Param("strategy_id")
	strategyID, err := strconv.Atoi(strategyIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID := uint(h.Repository.GetCreatorID())
	load, created, err := h.Repository.GetSystemLoadDraft(creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	err = h.Repository.AddStrategy(uint(strategyID), creatorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.SystemLoadID)
	status := http.StatusOK
	if created {
		ctx.Header("Location", fmt.Sprintf("/api/system_loads/%d", load.SystemLoadID))
		status = http.StatusCreated
	}
	ctx.JSON(status, serializer.SystemLoadToJSON(load, creatorLogin, moderatorLogin, completedCount))
}

func (h *Handler) DeleteFromSystemLoad(ctx *gin.Context) {
	strategyID, err := strconv.Atoi(ctx.Param("strategy_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	systemLoadID, err := strconv.Atoi(ctx.Param("system_load_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.DeleteStrategyFromSystemLoad(systemLoadID, strategyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.SystemLoadID)
	ctx.JSON(http.StatusOK, serializer.SystemLoadToJSON(load, creatorLogin, moderatorLogin, completedCount))
}

func (h *Handler) EditInSystemLoad(ctx *gin.Context) {
	strategyID, err := strconv.Atoi(ctx.Param("strategy_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	systemLoadID, err := strconv.Atoi(ctx.Param("system_load_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.SystemLoadStrategyJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	item, err := h.Repository.EditStrategyInSystemLoad(systemLoadID, strategyID, j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.SystemLoadStrategyToJSON(item))
}
