package controllers

import (
	"daily-diet-backend/models"
	"daily-diet-backend/repositories"
	"daily-diet-backend/services"
	"daily-diet-backend/utils/logger"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController interface {
	CreateUser(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
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
		authRouter.GET("/login", authController.SignIn)
		authRouter.GET("/user/:email", authController.GetUserByEmail)
	}
}

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

func (controller *authController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := controller.service.GetUserByEmail(ctx, email)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "error parsing request"})
	}
	ctx.JSON(http.StatusFound, gin.H{"user": user})
}

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
	ctx.Set("Authorization", "Bearer "+token)
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
