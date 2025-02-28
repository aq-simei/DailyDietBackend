package controllers

import (
	"daily-diet-backend/middlewares"
	"daily-diet-backend/repositories"
	"daily-diet-backend/services"
	"daily-diet-backend/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatsController interface {
	GetStats(ctx *gin.Context)
}

type userStatsController struct {
	service services.UserStatsService
}

func RegisterUserStatsRoutes(router *gin.RouterGroup, client *gorm.DB, authService services.AuthService) {
	userStatsRepo := repositories.NewUserStatsRepository(client)
	userStatsService := services.NewUserStatsService(userStatsRepo)
	userStatsController := NewUserStatsController(userStatsService)
	userStatsRouter := router.Group("/userstats")

	userStatsRouter.Use(middlewares.AuthMiddleware(authService))
	logger.Log(logger.DEBUG, "Registering auth routes")
	{
		userStatsRouter.GET("/find", userStatsController.GetStats)
	}
}
func NewUserStatsController(service services.UserStatsService) UserStatsController {
	return &userStatsController{service: service}
}

func RegisteredUserStatsRoutes(router *gin.RouterGroup, userStatsService services.UserStatsService) {
	userStatsController := NewUserStatsController(userStatsService)
	userStatsRouter := router.Group("/userstats")

	{
		userStatsRouter.GET("/stats", userStatsController.GetStats)
	}
}

func (controller *userStatsController) GetStats(ctx *gin.Context) {
	userId := ctx.Keys["userId"]
	parserId, err := uuid.Parse(userId.(string))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Error parsing userId"})
		return
	}
	stats, err := controller.service.GetStats(ctx, parserId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	ctx.JSON(200, stats)
}
