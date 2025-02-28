package services

import (
	"context"
	"daily-diet-backend/models"
	"daily-diet-backend/repositories"

	"github.com/google/uuid"
)

type UserStatsService interface {
	GetStats(c context.Context, userId uuid.UUID) (*models.UserStats, error)
}

type userStatsService struct {
	repo repositories.UserStatsRepository
}

func NewUserStatsService(repo repositories.UserStatsRepository) UserStatsService {
	return &userStatsService{repo: repo}
}

func (s *userStatsService) GetStats(c context.Context, userId uuid.UUID) (*models.UserStats, error) {
	return s.repo.GetStats(c, userId)
}
