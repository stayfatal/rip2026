package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"web_backend/internal/app/repository"
)

// Handler содержит зависимости HTTP-обработчиков.
type Handler struct {
	Repository *repository.Repository
}

// NewHandler создаёт новый Handler с переданным репозиторием.
func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// GetStrategies — обработчик главной страницы: список стратегий + иконка заявки.
func (h *Handler) GetStrategies(ctx *gin.Context) {
	var strategies []repository.ShardingStrategy
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		strategies, err = h.Repository.GetStrategies()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		strategies, err = h.Repository.GetStrategiesByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	systemLoads, err := h.Repository.GetSystemLoads()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"strategies":   strategies,
		"query":        searchQuery,
		"system_loads": systemLoads,
	})
}

// GetStrategy — обработчик страницы детальной информации о стратегии.
func (h *Handler) GetStrategy(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	strategy, err := h.Repository.GetStrategy(id)
	if err != nil {
		logrus.Error(err)
	}

	loadStrat, err := h.Repository.GetSystemLoadForStrategy(id)
	hasLoadInfo := err == nil && loadStrat != nil

	ctx.HTML(http.StatusOK, "strategy.html", gin.H{
		"strategy":    strategy,
		"loadStrat":   loadStrat,
		"hasLoadInfo": hasLoadInfo,
	})
}

// GetSystemLoad — обработчик страницы заявки (описание системы и нагрузка).
func (h *Handler) GetSystemLoad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	load, err := h.Repository.GetSystemLoad(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "system_load.html", gin.H{
		"load": load,
	})
}
