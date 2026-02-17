package handler

import (
	"fmt"
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

// GetStrategies — обработчик главной страницы: список стратегий + карточка заявки.
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

	calculations, err := h.Repository.GetCalculations()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"strategies":   strategies,
		"query":        searchQuery,
		"calculations": calculations,
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

	calcStrat, err := h.Repository.GetCalculationForStrategy(id)
	hasCalcInfo := err == nil && calcStrat != nil

	ctx.HTML(http.StatusOK, "strategy.html", gin.H{
		"strategy":    strategy,
		"calcStrat":   calcStrat,
		"hasCalcInfo": hasCalcInfo,
	})
}

// GetCalculation — обработчик страницы состава заявки (расчёта нагрузки).
func (h *Handler) GetCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	calc, err := h.Repository.GetCalculation(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "calculation.html", gin.H{
		"calc":         calc,
		"responseTime": fmt.Sprintf("%.1f", calc.ResultResponseTime),
	})
}
