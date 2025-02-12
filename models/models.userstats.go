package models

import "github.com/google/uuid"

type UserStats struct {
	ID              uuid.UUID `json:"id" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	UserID          uuid.UUID `json:"userId" gorm:"type:uuid"`
	RegisteredMeals int       `json:"registeredMeals" gorm:"default:0"`
	InDietMeals     int       `json:"inDietMeals" gorm:"default:0"`
	CurrentStreak   int       `json:"currentStreak" gorm:"default:0"`
	MaxStreak       int       `json:"maxStreak" gorm:"default:0"`
	User            *User     `json:"-" gorm:"foreignKey:UserID"`
}

func (UserStats) TableName() string {
	return "user_stats"
}

func (userStats *UserStats) UpdateMaxStreak(
	inDiet bool,
	currentMax int,
	current int,
) int {
	if inDiet {
		if (current + 1) > currentMax {
			return current + 1
		}
		if currentMax == 0 {
			return 1
		}
	}
	if currentMax == 0 {
		return 0
	}
	return userStats.MaxStreak
}

func (userStats *UserStats) UpdateCurrentStreak(
	inDiet bool,
	current int,
) int {
	if inDiet {
		return current + 1
	}
	return 0
}

func (userStats *UserStats) UpdateInDietMeals(
	inDiet bool,
	isNewMeal bool,
	isDelete bool,
	wasInDiet bool,
	current int,
) int {
	if isNewMeal {
		if inDiet {
			return current + 1
		}
		return current
	}

	if isDelete {
		if wasInDiet {
			return current - 1
		}
		return current
	}

	// Handling updates
	if inDiet && !wasInDiet {
		return current + 1
	}
	if !inDiet && wasInDiet {
		return current - 1
	}
	return current
}

func (userStats *UserStats) UpdateRegisteredMeals(isNew bool, current int) int {
	if isNew {
		return current + 1
	}
	return current - 1
}
