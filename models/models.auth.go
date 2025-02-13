package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginDTO struct {
	Email    string  `json:"email" binding:"required"`
	Password string  `json:"password" binding:"required"`
	DeviceID *string `json:"device_id"`
}

type UserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtTokenClaims struct {
	Email string    `json:"email"`
	ID    uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}
