package handler

import (
	"net/http"
	"strconv"
	"web_backend/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetStrategies(ctx *gin.Context) {
	var strategies []ds.ShardingStrategy
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		strategies, err = h.Repository.GetStrategies()
	} else {
		strategies, err = h.Repository.GetStrategiesByTitle(searchQuery)
	}
	if err != nil {
		logrus.Error(err)
	}

	creatorID := uint(1)
	loadCount := h.Repository.GetSystemLoadStrategyCount(creatorID)
	activeLoadID := h.Repository.GetActiveSystemLoadID(creatorID)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"strategies":      strategies,
		"query":           searchQuery,
		"system_load_count": loadCount,
		"system_load_id":   activeLoadID,
		"minioUrl":         h.Config.MinioURL,
	})
}

func (h *Handler) GetStrategy(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	strategy, err := h.Repository.GetStrategy(id)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "strategy.html", gin.H{
		"strategy": strategy,
		"minioUrl": h.Config.MinioURL,
	})
}
