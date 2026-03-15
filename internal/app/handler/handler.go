package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"web_backend/internal/app/config"
	"web_backend/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
}

func NewHandler(r *repository.Repository, cfg *config.Config) *Handler {
	return &Handler{
		Repository: r,
		Config:     cfg,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")

	strategies := api.Group("/strategies")
	{
		strategies.GET("", h.GetStrategies)
		strategies.GET("/:id", h.GetStrategy)
		strategies.POST("", h.CreateStrategy)
	}

	systemLoads := api.Group("/system_loads")
	{
		systemLoads.GET("/cart", h.GetSystemLoadCart)
		systemLoads.GET("", h.GetAllSystemLoads)
		systemLoads.GET("/:id", h.GetSystemLoad)
		systemLoads.PUT("/:id", h.EditSystemLoad)
		systemLoads.PUT("/:id/form", h.FormSystemLoad)
		systemLoads.PUT("/:id/finish", h.FinishSystemLoad)
		systemLoads.DELETE("/:id", h.DeleteSystemLoad)
	}

	mm := api.Group("/system_load_strategies")
	{
		mm.POST("/add/:strategy_id", h.AddToSystemLoad)
		mm.DELETE("/:strategy_id/:system_load_id", h.DeleteFromSystemLoad)
		mm.PUT("/:strategy_id/:system_load_id", h.EditInSystemLoad)
	}

	users := api.Group("/users")
	{
		users.POST("/register", h.CreateUser)
		users.POST("/login", h.SignIn)
		users.POST("/logout", h.SignOut)
	}
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}
