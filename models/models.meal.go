package models

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID          uuid.UUID `json:"id" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Date        time.Time `json:"date" gorm:"not null"`
	Time        time.Time `json:"time" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	InDiet      bool      `json:"in_diet" gorm:"not null"`
}

func (Meal) TableName() string {
	return "meals"
}

type EditMealDTO struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	Time        *time.Time `json:"time,omitempty"`
	InDiet      *bool      `json:"in_diet,omitempty"`
}

type CreateMealDTO struct {
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description,omitempty"`
	Date        time.Time `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Time        time.Time `json:"time" binding:"required"` // Format: HH:mm
	InDiet      bool      `json:"in_diet" binding:"boolean"`
}

type GetMealDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Time        time.Time `json:"time"`
	InDiet      bool      `json:"in_diet"`
}
