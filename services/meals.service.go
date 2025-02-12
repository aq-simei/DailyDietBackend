package services

import (
	"context"
	"daily-diet-backend/models"
	"daily-diet-backend/repositories"

	"github.com/google/uuid"
)

type MealsService interface {
	GetMeals(c context.Context, userId uuid.UUID) ([]models.Meal, error)
	CreateMeal(c context.Context, data models.CreateMealDTO, userId uuid.UUID) (*models.Meal, error)
	DeleteMeal(c context.Context, mealId string, userId uuid.UUID) error
	EditMeal(c context.Context, mealId string, userId uuid.UUID, data models.EditMealDTO) (*models.Meal, error)
}

type mealsService struct {
	repo repositories.MealsRepository
}

func NewMealsService(repo repositories.MealsRepository) MealsService {
	return &mealsService{repo: repo}
}

func (service *mealsService) GetMeals(c context.Context, userId uuid.UUID) ([]models.Meal, error) {
	return service.repo.GetMeals(c, userId)
}

func (service *mealsService) CreateMeal(c context.Context, data models.CreateMealDTO, userId uuid.UUID) (*models.Meal, error) {
	return service.repo.CreateMeal(c, data, userId)
}

func (service *mealsService) DeleteMeal(c context.Context, mealId string, userId uuid.UUID) error {
	return service.repo.DeleteMeal(c, mealId, userId)
}

func (service *mealsService) EditMeal(
	c context.Context,
	mealId string,
	userId uuid.UUID,
	data models.EditMealDTO,
) (*models.Meal, error) {
	return service.repo.EditMeal(c, mealId, userId, data)
}
