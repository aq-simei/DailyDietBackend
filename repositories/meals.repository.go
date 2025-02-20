package repositories

import (
	"context"

	"daily-diet-backend/models"
	"daily-diet-backend/utils/errors"
	"daily-diet-backend/utils/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MealsRepository interface {
	GetMeals(c context.Context, userId uuid.UUID) ([]models.Meal, error)
	CreateMeal(c context.Context, data models.CreateMealDTO, userId uuid.UUID) (*models.Meal, error)
	DeleteMeal(c context.Context, mealId string, userId uuid.UUID) error
	EditMeal(c context.Context, mealId string, userId uuid.UUID, data models.EditMealDTO) (*models.Meal, error)
	GetMeal(c context.Context, mealId string, userId uuid.UUID) (*models.Meal, error)
}

type mealsRepository struct {
	database *gorm.DB
}

func NewMealsRepository(client *gorm.DB) MealsRepository {
	return &mealsRepository{database: client}
}

func (repo *mealsRepository) GetMeals(
	c context.Context,
	userId uuid.UUID,
) ([]models.Meal, error) {
	var meals []models.Meal
	if err := repo.database.WithContext(c).Where("user_id = ?", userId).Find(&meals).Error; err != nil {
		return nil, err
	}
	return meals, nil
}

func (repo *mealsRepository) CreateMeal(
	c context.Context,
	data models.CreateMealDTO,
	userId uuid.UUID) (*models.Meal, error) {
	var meal *models.Meal
	var txErr error

	txErr = repo.database.Transaction(
		func(tx *gorm.DB) error {
			meal = &models.Meal{
				Name:   data.Name,
				UserID: userId,
				Date:   data.Date,
				Time:   data.Time,
				InDiet: data.InDiet,
			}

			if data.Description != nil {
				meal.Description = *data.Description
			}

			if err := tx.Create(meal).Error; err != nil {
				txErr = err
				return err // rollback
			}

			if err := repo.handlePostCreate(tx, c, data, userId); err != nil {
				txErr = err
				return err // rollback
			}

			return nil // transaction commit
		})
	if txErr != nil {
		logger.Log(logger.ERROR, "Error creating meal")
		logger.Log(logger.ERROR, txErr.Error())
		return nil, txErr
	}
	return meal, nil
}

func (repo *mealsRepository) handlePostCreate(
	tx *gorm.DB,
	c context.Context,
	data models.CreateMealDTO,
	userId uuid.UUID,
) error {
	var existingUser models.User
	if err := tx.WithContext(c).First(&existingUser, userId).Error; err != nil {
		logger.Log(logger.ERROR, "Error finding user: "+err.Error())
		return err
	}

	// Find or create stats
	var stats models.UserStats
	result := tx.WithContext(c).Where("user_id = ?", userId).First(&stats)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			stats = models.UserStats{
				UserID: userId,
				CurrentStreak: stats.UpdateCurrentStreak(
					data.InDiet, 0,
				),
				MaxStreak: stats.UpdateMaxStreak(
					data.InDiet,
					0,
					0,
				),
				InDietMeals: stats.UpdateInDietMeals(
					data.InDiet,
					true,
					false,
					false,
					0,
				),
				RegisteredMeals: 1,
			}
			if err := tx.Create(&stats).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	} else {
		// Update stats
		stats.InDietMeals = stats.UpdateInDietMeals(
			data.InDiet,
			true,
			false,
			false,
			stats.InDietMeals,
		)
		stats.RegisteredMeals = stats.UpdateRegisteredMeals(
			true,
			stats.RegisteredMeals,
		)
		stats.MaxStreak = stats.UpdateMaxStreak(
			data.InDiet,
			stats.MaxStreak,
			stats.CurrentStreak,
		)
		stats.CurrentStreak = stats.UpdateCurrentStreak(
			data.InDiet,
			stats.CurrentStreak,
		)
		if err := tx.Save(&stats).Error; err != nil {
			return err
		}
	}

	// Update associations
	return tx.Save(&existingUser).Error
}

func (repo *mealsRepository) DeleteMeal(
	c context.Context,
	mealId string,
	userId uuid.UUID,
) error {
	var toDeleteMeal models.Meal
	if err := repo.database.WithContext(c).
		Where("id = ?", mealId).
		First(&toDeleteMeal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewError(
				errors.NotFound,
				"no meal with id -> "+mealId, err,
			)
		}
		return errors.NewError(
			errors.Internal,
			"could not find meal with id ->"+mealId,
			err,
		)
	}

	var existingStats models.UserStats
	err := repo.database.Where("user_id = ?", userId).First(&existingStats).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Log(logger.ERROR, "Error finding user stats")
			logger.Log(logger.ERROR, err.Error())
			return err
		}
		logger.Log(logger.INFO, "No stats found, initializing with zero values")
	} else {
	}

	if err := repo.database.WithContext(c).Delete(&toDeleteMeal).Error; err != nil {
		return err
	}

	// update stats
	existingStats.InDietMeals = existingStats.UpdateInDietMeals(
		toDeleteMeal.InDiet,
		false,
		true,
		toDeleteMeal.InDiet,
		existingStats.InDietMeals,
	)
	existingStats.RegisteredMeals = existingStats.UpdateRegisteredMeals(
		false,
		existingStats.RegisteredMeals,
	)
	existingStats.MaxStreak = existingStats.UpdateMaxStreak(
		toDeleteMeal.InDiet,
		existingStats.MaxStreak,
		existingStats.CurrentStreak,
	)
	existingStats.CurrentStreak = existingStats.UpdateCurrentStreak(
		toDeleteMeal.InDiet,
		existingStats.CurrentStreak,
	)
	if err := repo.database.Save(&existingStats).Error; err != nil {
		return errors.NewError(
			errors.Internal,
			"error updating user stats on delete -> ",
			err,
		)
	}

	return nil
}

func (repo *mealsRepository) EditMeal(
	c context.Context,
	mealId string,
	userId uuid.UUID,
	data models.EditMealDTO,
) (*models.Meal, error) {
	var toEditMeal *models.Meal

	err := repo.database.WithContext(c).
		Where("id = ?", mealId).
		First(&toEditMeal).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewError(
				errors.NotFound,
				"no meal with id -> "+mealId,
				err,
			)
		}
		return nil, errors.NewError(
			errors.Internal,
			"could not find meal with id ->"+mealId,
			err,
		)
	}

	if toEditMeal.UserID != userId {
		return nil, errors.NewError(
			errors.Forbidden,
			"you are not allowed to edit this meal",
			nil,
		)
	}

	if data.Name != nil {
		toEditMeal.Name = *data.Name
	}
	if data.Description != nil {
		toEditMeal.Description = *data.Description
	}
	if data.Date != nil {
		toEditMeal.Date = *data.Date
	}
	if data.Time != nil {
		toEditMeal.Time = *data.Time
	}
	if data.InDiet != nil {
		toEditMeal.InDiet = *data.InDiet
	}

	if err := repo.database.
		WithContext(c).
		Save(&toEditMeal).Error; err != nil {
		return nil, errors.NewError(
			errors.Internal,
			"error updating meal",
			err,
		)
	}

	return toEditMeal, nil
}

func (repo *mealsRepository) GetMeal(
	c context.Context,
	mealId string,
	userId uuid.UUID,
) (*models.Meal, error) {
	var meal *models.Meal
	if err := repo.database.WithContext(c).
		Where("id = ? AND user_id = ?", mealId, userId).
		First(&meal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewError(
				errors.NotFound,
				"no meal with id -> "+mealId,
				err,
			)
		}
		return nil, errors.NewError(
			errors.Internal,
			"could not find meal with id ->"+mealId,
			err,
		)
	}
	return meal, nil
}
