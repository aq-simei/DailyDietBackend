package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string    `gorm:"type:varchar(255);primaryKey;not null" json:"token"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	DeviceID  *string   `gorm:"type:varchar(255);" json:"device_id"`
	ExpireAt  time.Time `gorm:"type:timestamp;not null" json:"expire_at"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;" json:"updated_at"`
	Revoked   bool      `gorm:"type:boolean;not null;default:false;index" json:"revoked"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type CreateRefreshTokenDTO struct {
	UserID   uuid.UUID `json:"user_id"`
	DeviceID *string   `json:"device_id"`
}

type ValidateRefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

// TableName specifies the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
