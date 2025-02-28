package repositories

import (
	"context"
	"daily-diet-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatsRepository interface {
	GetStats(c context.Context, userId uuid.UUID) (*models.UserStats, error)
}

type userStatsRepository struct {
	database *gorm.DB
}

func NewUserStatsRepository(client *gorm.DB) UserStatsRepository {
	return &userStatsRepository{database: client}
}

func (repo *userStatsRepository) GetStats(
	c context.Context,
	userId uuid.UUID,
) (*models.UserStats, error) {
	var stats models.UserStats
	if err := repo.database.WithContext(c).Where("user_id = ?", userId).Find(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}
