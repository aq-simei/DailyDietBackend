package services

import (
	"context"
	"daily-diet-backend/models"
	"daily-diet-backend/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	CreateUser(c context.Context, data models.CreateUserDTO) (*models.User, error)
	GetUserByEmail(c context.Context, email string) (*models.User, error)
	Login(c context.Context, data models.LoginDTO) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (*models.JwtTokenClaims, error)
	ValidateRefreshToken(c context.Context, tokenString string) (*models.RefreshToken, error)
}

type authService struct {
	Repo      repositories.UserRepository
	JwtSecret []byte
}

func NewAuthService(repo repositories.UserRepository, jwtSecret []byte) AuthService {
	return &authService{Repo: repo, JwtSecret: jwtSecret}
}

func (service *authService) CreateUser(c context.Context, data models.CreateUserDTO) (*models.User, error) {
	return service.Repo.CreateUser(c, data)
}

func (service *authService) GetUserByEmail(c context.Context, email string) (*models.User, error) {
	return service.Repo.GetUserByEmail(c, email)
}

func (service *authService) Login(c context.Context, data models.LoginDTO) (*models.LoginResponse, error) {
	userLogin, err := service.Repo.Login(c, data)
	if err != nil {
		return nil, err
	}
	/*
		Generate JWT token -> 1 hour expiration
	*/
	claims := &models.JwtTokenClaims{
		Email: data.Email,
		/* store userId, better for fetches latter */
		ID: userLogin.User.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "daily-diet-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(service.JwtSecret)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{
		Token:        signedToken,
		RefreshToken: userLogin.RefreshToken,
		User:         userLogin.User,
	}, nil
}

func (service *authService) ValidateToken(tokenString string) (*models.JwtTokenClaims, error) {
	claims := &models.JwtTokenClaims{} // declare new empty claims
	// parse jwt with base JwtTokenClaims structure
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return service.JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JwtTokenClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (service *authService) ValidateRefreshToken(c context.Context, tokenString string) (*models.RefreshToken, error) {
	userId := c.Value("userId").(string)
	return service.Repo.ValidateRefreshToken(c, tokenString, userId)
}
