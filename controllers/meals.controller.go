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
	GetMeal(ctx *gin.Context)
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
		mealsRouter.GET("/:mealId", mealsController.GetMeal)
	}
}

// CreateMeal godoc
// @Summary Create a new meal
// @Description Creates a new meal for the authenticated user
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param meal body models.CreateMealDTO true "Meal details"
// @Success 201 {object} models.Meal
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meals/new [post]
func (controller *mealsController) CreateMeal(ctx *gin.Context) {
	userId := ctx.Keys["userId"].(string)
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not parse userId"})
		return
	}
	var req models.CreateMealDTO
	logger.Log(logger.DEBUG, "Creating meal")

	if err := ctx.BindJSON(&req); err != nil {
		logger.Log(logger.ERROR, err.Error())
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

// EditMeal godoc
// @Summary Edit an existing meal
// @Description Modifies an existing meal for the authenticated user
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mealId path string true "Meal ID"
// @Param meal body models.EditMealDTO true "Updated meal details"
// @Success 200 {object} models.Meal
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meals/edit/{mealId} [patch]
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
		logger.Log(logger.ERROR, err.Error())
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

// GetMeals godoc
// @Summary List all meals
// @Description Retrieves all meals for the authenticated user
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Meal
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meals/list [get]
func (controller *mealsController) GetMeals(ctx *gin.Context) {
	userId := ctx.Keys["userId"].(string)
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not parse userId"})
		return
	}
	if userId == "" {
		ctx.JSON(404, gin.H{"error": "userId not found"})
		return
	}

	meals, err := controller.service.GetMeals(ctx, parsedUserId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, meals)
}

// DeleteMeal godoc
// @Summary Delete a meal
// @Description Deletes a specific meal for the authenticated user
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mealId path string true "Meal ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meals/delete/{mealId} [delete]
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

func (controller *mealsController) GetMeal(ctx *gin.Context) {
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

	meal, err := controller.service.GetMeal(ctx, mealId, parsedUserId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, meal)
}
