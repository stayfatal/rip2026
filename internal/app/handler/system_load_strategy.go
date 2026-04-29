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

// AddToSystemLoad godoc
// @Summary Добавить стратегию в заявку
// @Description Добавляет стратегию в заявку-черновик пользователя
// @Tags system_load_strategies
// @Produce json
// @Param strategy_id path int true "ID стратегии"
// @Success 200 {object} serializer.SystemLoadJSON "Заявка с добавленной стратегией"
// @Success 201 {object} serializer.SystemLoadJSON "Создана новая заявка"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Стратегия не найдена"
// @Failure 409 {object} map[string]string "Стратегия уже в заявке"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /system_load_strategies/add/{strategy_id} [post]
func (h *Handler) AddToSystemLoad(ctx *gin.Context) {
	strategyIDStr := ctx.Param("strategy_id")
	strategyID, err := strconv.Atoi(strategyIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
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

// DeleteFromSystemLoad godoc
// @Summary Удалить стратегию из заявки
// @Description Удаляет связь стратегии и заявки (только черновик)
// @Tags system_load_strategies
// @Produce json
// @Param strategy_id path int true "ID стратегии"
// @Param system_load_id path int true "ID заявки"
// @Success 200 {object} serializer.SystemLoadJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /system_load_strategies/{strategy_id}/{system_load_id} [delete]
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

// EditInSystemLoad godoc
// @Summary Изменить данные стратегии в заявке
// @Description Обновляет параметры стратегии в конкретной заявке (только черновик)
// @Tags system_load_strategies
// @Accept json
// @Produce json
// @Param strategy_id path int true "ID стратегии"
// @Param system_load_id path int true "ID заявки"
// @Param data body serializer.SystemLoadStrategyJSON true "Новые данные"
// @Success 200 {object} serializer.SystemLoadStrategyJSON "Обновленные данные"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /system_load_strategies/{strategy_id}/{system_load_id} [put]
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
