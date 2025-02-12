package router

import (
	"daily-diet-backend/controllers"
	"daily-diet-backend/repositories"
	"daily-diet-backend/services"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func NewRouter(client *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Swagger setup
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := router.Group("/v1")
	jwt := []byte(os.Getenv("JWT_SECRET"))
	usersRepo := repositories.NewUserRepository(client)
	authService := services.NewAuthService(usersRepo, jwt)

	controllers.RegisterAuthRoutes(v1, client)
	controllers.RegisteredMealsRoutes(v1, client, authService)

	return router
}
