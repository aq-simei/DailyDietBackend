package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// Primary key field using UUID
	ID uuid.UUID `json:"id" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	// Unique email field that cannot be null
	Email string `json:"email" gorm:"unique;not null"`
	// Required name field
	Name string `json:"name" gorm:"not null"`
	// Required password field
	Password string `json:"password" gorm:"not null"`
	// Automatically managed timestamp fields
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// One-to-Many relation with Meal
	// - foreignKey:UserID: specifies the foreign key field in Meal table
	// - constraint:OnDelete:CASCADE: deletes related meals when user is deleted
	Meals []Meal `json:"meals" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	// Change the type from uuid.UUID to UserStats
	UserStats UserStats `json:"user_stats" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	RefreshToken *RefreshToken `json:"refresh_token,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "users"
}
