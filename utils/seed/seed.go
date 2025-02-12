package seed

import (
	"context"
	"time"

	"daily-diet-backend/models"
	"daily-diet-backend/utils/crypt"
	"daily-diet-backend/utils/errors"
	"daily-diet-backend/utils/logger"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB, ctx context.Context) error {
	// check if there is user in database
	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Log(logger.INFO, "No users found in database, seeding...")
		} else {
			logger.Log(logger.ERROR, "Error finding users :: "+result.Error.Error())
			return errors.NewError(errors.Internal, "Error finding users :: ", result.Error)
		}
	}

	if len(users) == 0 {
		// Create test user
		hashedPassword, _ := crypt.HashPassword("test122")

		toCreateUser := &models.User{
			Email:    "leo@mail.com",
			Name:     "Leo Messi",
			Password: string(hashedPassword),
		}
		meals := []models.Meal{
			{
				Name:        "Healthy Breakfast",
				Date:        time.Date(2024, 1, 15, 8, 30, 0, 0, time.Local),
				Time:        time.Date(2024, 1, 15, 8, 30, 0, 0, time.Local),
				InDiet:      true,
				Description: "Oatmeal with fruits and honey",
				UserID:      toCreateUser.ID,
			},
			{
				Name:        "Fast Food Lunch",
				Date:        time.Date(2024, 1, 15, 12, 45, 0, 0, time.Local),
				Time:        time.Date(2024, 1, 15, 12, 45, 0, 0, time.Local),
				InDiet:      false,
				Description: "Double cheeseburger with fries",
				UserID:      toCreateUser.ID,
			},
			{
				Name:        "Healthy Dinner",
				Date:        time.Date(2024, 1, 15, 19, 0, 0, 0, time.Local),
				Time:        time.Date(2024, 1, 15, 19, 0, 0, 0, time.Local),
				InDiet:      true,
				Description: "Grilled chicken with salad",
				UserID:      toCreateUser.ID,
			},
			{
				Name:        "Late Night Snack",
				Date:        time.Date(2024, 1, 15, 23, 15, 0, 0, time.Local),
				Time:        time.Date(2024, 1, 15, 23, 15, 0, 0, time.Local),
				InDiet:      false,
				Description: "Chocolate cake and ice cream",
				UserID:      toCreateUser.ID,
			},
		}
		result := db.WithContext(ctx).Create(toCreateUser)
		if result.Error != nil {
			err := result.Error
			logger.Log(logger.ERROR, "Error creating user :: "+err.Error())
			return errors.NewError(errors.Internal, "Error creating user :: ", err)
		}

		// Create meals

		for _, meal := range meals {
			result = db.WithContext(ctx).Create(&meal)
			if result.Error != nil {
				err := result.Error
				logger.Log(logger.ERROR, "Error creating meal :: "+err.Error())
				return errors.NewError(errors.Internal, "Error creating meal :: ", err)
			}
		}

		toCreateUserStats := &models.UserStats{
			UserID:          toCreateUser.ID,
			RegisteredMeals: 4,
			InDietMeals:     2,
			CurrentStreak:   0, // Reset to 0 since last meal was not in diet
			MaxStreak:       2, // Had 2 meals in diet in sequence
		}

		result = db.WithContext(ctx).Create(toCreateUserStats)
		if result.Error != nil {
			err := result.Error
			logger.Log(logger.ERROR, "Error creating user stats :: "+err.Error())
			return errors.NewError(errors.Internal, "Error creating user stats :: ", err)
		}
		logger.Log(logger.INFO, "Seed completed successfully")
	} else {
		logger.Log(logger.WARNING, "Database already seeded, check utils/seed/seed.go")
		return nil
	}

	return nil
}
