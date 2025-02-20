package controllers

import (
	"daily-diet-backend/models"
	"daily-diet-backend/repositories"
	"daily-diet-backend/services"
	"daily-diet-backend/utils/logger"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthController interface {
	CreateUser(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
	RefreshTokenLogin(ctx *gin.Context)
}

type authController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) AuthController {
	return &authController{service: service}
}

func RegisterAuthRoutes(router *gin.RouterGroup, client *gorm.DB) {
	jwt := []byte(os.Getenv("JWT_SECRET"))
	usersRepo := repositories.NewUserRepository(client)
	authService := services.NewAuthService(usersRepo, jwt)
	authController := NewAuthController(authService)

	authRouter := router.Group("/auth")
	logger.Log(logger.DEBUG, "Registering auth routes")
	{
		authRouter.POST("/register", authController.CreateUser)
		authRouter.POST("/login", authController.SignIn)
		authRouter.POST("/login/token", authController.RefreshTokenLogin)
		authRouter.GET("/user/:email", authController.GetUserByEmail)
	}
}

// CreateUser godoc
// @Summary Register new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.CreateUserDTO true "User registration details"
// @Success 201 {object} models.UserDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (controller *authController) CreateUser(ctx *gin.Context) {
	var req models.CreateUserDTO
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "error parsing request"})
		return
	}

	user, err := controller.service.CreateUser(ctx, req)
	// is error from NewError
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var createdUser models.UserDTO = models.UserDTO{
		Email: user.Email,
		Name:  user.Name,
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Retrieves user information by email address
// @Tags auth
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 302 {object} map[string]models.User
// @Failure 400 {object} map[string]string
// @Router /auth/user/{email} [get]
func (controller *authController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := controller.service.GetUserByEmail(ctx, email)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "error parsing request"})
	}
	ctx.JSON(http.StatusFound, gin.H{"user": user})
}

type SuccessResponse struct {
	Token string `json:"token"`
}

// SignIn godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginDTO true "Login credentials"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [get]
func (controller *authController) SignIn(ctx *gin.Context) {
	var req models.LoginDTO
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "error parsing request"})
		return
	}

	token, err := controller.service.Login(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Set("Authorization", "Bearer "+token.Token)
	ctx.JSON(http.StatusOK,
		gin.H{
			"token":         token.Token,
			"refresh_token": token.RefreshToken,
		})
}

func (controller *authController) RefreshTokenLogin(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "error parsing request"})
		return
	}

	validateRefreshTokenResponse, err := controller.service.ValidateRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	updatedRefreshToken, err := controller.service.UpdateRefreshToken(ctx, *validateRefreshTokenResponse.RefreshToken, validateRefreshTokenResponse.UserID.String())

	if err != nil {
		ctx.JSON(500, gin.H{"error": "error updating refresh token ::" + err.Error()})
		return
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	/*
		Generate JWT token -> 1 hour expiration
	*/
	claims := &models.JwtTokenClaims{
		Email: *validateRefreshTokenResponse.UserEmail,
		/* store userId, better for fetches latter */
		ID: *validateRefreshTokenResponse.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "daily-diet-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		ctx.JSON(500, gin.H{"Jwt Sign Error": err.Error()})
		return
	}

	ctx.Set("Authorization", "Bearer "+signedToken)
	ctx.JSON(http.StatusOK, gin.H{
		"refresh_token": updatedRefreshToken.Token,
		"jwt_token":     signedToken,
		"user_email":    validateRefreshTokenResponse.UserEmail,
		"user_id":       updatedRefreshToken.UserID,
	})
}
