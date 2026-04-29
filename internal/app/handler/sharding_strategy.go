package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/ds"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

// GetStrategies godoc
// @Summary Получить список стратегий
// @Description Возвращает все стратегии или фильтрует по названию
// @Tags strategies
// @Produce json
// @Param Title query string false "Название стратегии для поиска"
// @Success 200 {array} serializer.ShardingStrategyJSON "Список стратегий"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /strategies [get]
func (h *Handler) GetStrategies(ctx *gin.Context) {
	var strategies []ds.ShardingStrategy
	var err error
	searchQuery := ctx.Query("Title")
	if searchQuery == "" {
		strategies, err = h.Repository.GetStrategies()
	} else {
		strategies, err = h.Repository.GetStrategiesByTitle(searchQuery)
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.ShardingStrategyJSON, 0, len(strategies))
	for _, s := range strategies {
		resp = append(resp, serializer.ShardingStrategyToJSON(s))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetStrategy godoc
// @Summary Получить стратегию по ID
// @Description Возвращает информацию о стратегии по идентификатору
// @Tags strategies
// @Produce json
// @Param id path int true "ID стратегии"
// @Success 200 {object} serializer.ShardingStrategyJSON "Данные стратегии"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Стратегия не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /strategies/{id} [get]
func (h *Handler) GetStrategy(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	strategy, err := h.Repository.GetStrategy(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.ShardingStrategyToJSON(*strategy))
}

// CreateStrategy godoc
// @Summary Создать стратегию
// @Description Создает новую стратегию шардинга (с возможной загрузкой изображения/видео)
// @Tags strategies
// @Accept json
// @Produce json
// @Param strategy body serializer.ShardingStrategyJSON true "Данные новой стратегии"
// @Success 201 {object} serializer.ShardingStrategyJSON "Созданная стратегия"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /strategies [post]
func (h *Handler) CreateStrategy(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")
	var j serializer.ShardingStrategyJSON
	if strings.HasPrefix(contentType, "application/json") {
		if err := ctx.BindJSON(&j); err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	} else {
		title := ctx.PostForm("title")
		desc := ctx.PostForm("description")
		if title == "" || desc == "" {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("title and description are required"))
			return
		}
		latency := 1.0
		throughput := 1.0
		reliability := 0.9
		if v := ctx.PostForm("latency_coefficient"); v != "" {
			fmt.Sscanf(v, "%f", &latency)
		}
		if v := ctx.PostForm("throughput_coefficient"); v != "" {
			fmt.Sscanf(v, "%f", &throughput)
		}
		if v := ctx.PostForm("reliability_coefficient"); v != "" {
			fmt.Sscanf(v, "%f", &reliability)
		}
		j = serializer.ShardingStrategyJSON{
			Title:                  title,
			Description:            desc,
			LatencyCoefficient:     latency,
			ThroughputCoefficient:  throughput,
			ReliabilityCoefficient: reliability,
		}
	}

	strategy, err := h.Repository.CreateStrategy(j)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if imageFile, err := ctx.FormFile("image"); err == nil {
		s, err := h.Repository.AddPhoto(ctx, int(strategy.StrategyID), imageFile)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		strategy = *s
	}
	if videoFile, err := ctx.FormFile("video"); err == nil {
		s, err := h.Repository.AddVideo(ctx, int(strategy.StrategyID), videoFile)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		strategy = *s
	}

	ctx.Header("Location", fmt.Sprintf("/api/strategies/%d", strategy.StrategyID))
	ctx.JSON(http.StatusCreated, serializer.ShardingStrategyToJSON(strategy))
}
