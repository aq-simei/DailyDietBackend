package controllers

import (
	"daily-diet-backend/middlewares"
	"daily-diet-backend/models"
	"daily-diet-backend/repositories"
	"daily-diet-backend/services"
	"daily-diet-backend/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MealsController interface {
	CreateMeal(ctx *gin.Context)
	EditMeal(ctx *gin.Context)
	GetMeals(ctx *gin.Context)
	DeleteMeal(ctx *gin.Context)
}

type mealsController struct {
	service services.MealsService
}

func NewMealsController(service services.MealsService) MealsController {
	return &mealsController{service: service}
}

func RegisteredMealsRoutes(router *gin.RouterGroup, client *gorm.DB, authService services.AuthService) {
	mealsRepo := repositories.NewMealsRepository(client)
	mealsService := services.NewMealsService(mealsRepo)
	mealsController := NewMealsController(mealsService)
	mealsRouter := router.Group("/meals")

	mealsRouter.Use(middlewares.AuthMiddleware(authService))
	logger.Log(logger.DEBUG, "Registering auth routes")
	{
		mealsRouter.POST("/new", mealsController.CreateMeal)
		mealsRouter.GET("/list", mealsController.GetMeals)
		mealsRouter.PATCH("edit/:mealId", mealsController.EditMeal)
		mealsRouter.DELETE("delete/:mealId", mealsController.DeleteMeal)
	}
}

func (controller *mealsController) CreateMeal(ctx *gin.Context) {
	userId := ctx.Keys["userId"].(string)
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not parse userId"})
		return
	}
	var req models.CreateMealDTO

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	meal, err := controller.service.CreateMeal(ctx, req, parsedUserId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, meal)
}
func (controller *mealsController) EditMeal(ctx *gin.Context) {
	userId := ctx.Keys["userId"].(string)
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not parse userId"})
		return
	}
	mealId := ctx.Param("mealId")
	if mealId == "" {
		ctx.JSON(400, gin.H{"error": "mealId not found"})
		return
	}
	var req models.EditMealDTO
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	meal, err := controller.service.EditMeal(ctx, mealId, parsedUserId, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, meal)
}

func (controller *mealsController) GetMeals(ctx *gin.Context) {
	userId := ctx.Keys["userId"].(string)
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not parse userId"})
		return
	}
	if userId == "" {
		ctx.JSON(400, gin.H{"error": "userId not found"})
		return
	}

	meals, err := controller.service.GetMeals(ctx, parsedUserId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, meals)
}

func (controller *mealsController) DeleteMeal(ctx *gin.Context) {
	userId := ctx.Keys["userId"].(string)
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not parse userId"})
		return
	}
	mealId := ctx.Param("mealId")
	if mealId == "" {
		ctx.JSON(400, gin.H{"error": "mealId not found"})
		return
	}

	err = controller.service.DeleteMeal(ctx, mealId, parsedUserId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(204, nil)
}
