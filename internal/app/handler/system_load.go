package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

func (h *Handler) GetSystemLoadCart(ctx *gin.Context) {
	creatorID := uint(h.Repository.GetCreatorID())
	count := h.Repository.GetSystemLoadStrategyCount(creatorID)
	if count == 0 {
		load, err := h.Repository.CheckCurrentDraft(creatorID)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":          "no_draft",
				"strategies_count": 0,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"id":               load.SystemLoadID,
			"strategies_count": 0,
		})
		return
	}
	loadID := h.Repository.GetActiveSystemLoadID(creatorID)
	ctx.JSON(http.StatusOK, gin.H{
		"id":               loadID,
		"strategies_count": count,
	})
}

func (h *Handler) GetAllSystemLoads(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from, to time.Time
	if fromDate != "" {
		t, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = t
	}
	toDate := ctx.Query("to-date")
	if toDate != "" {
		t, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = t
	}
	status := ctx.Query("status")
	loads, err := h.Repository.GetAllSystemLoads(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.SystemLoadJSON, 0, len(loads))
	for _, load := range loads {
		creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
		completedCount, _ := h.Repository.GetCompletedItemCount(load.SystemLoadID)
		resp = append(resp, serializer.SystemLoadToJSON(load, creatorLogin, moderatorLogin, completedCount))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSystemLoad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.GetSingleSystemLoad(id)
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
	items, err := h.Repository.GetSystemLoadItems(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.SystemLoadID)
	itemsResp := make([]serializer.SystemLoadStrategyDetailJSON, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, serializer.SystemLoadStrategyDetailToJSON(item))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"system_load": serializer.SystemLoadToJSON(load, creatorLogin, moderatorLogin, completedCount),
		"strategies":  itemsResp,
	})
}

func (h *Handler) EditSystemLoad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.SystemLoadJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.EditSystemLoad(id, j)
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

func (h *Handler) FormSystemLoad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.FormSystemLoad(id)
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

func (h *Handler) FinishSystemLoad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var statusJSON serializer.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.FinishSystemLoad(id, statusJSON.Status)
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

func (h *Handler) DeleteSystemLoad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = h.Repository.DeleteSystemLoad(id)
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
	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка удалена"})
}
