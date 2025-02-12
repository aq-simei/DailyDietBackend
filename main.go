package main

import (
	"context"
	"daily-diet-backend/database"
	"daily-diet-backend/models"
	"daily-diet-backend/router"
	"daily-diet-backend/utils/logger"
	"daily-diet-backend/utils/seed"
	"daily-diet-backend/utils/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func main() {
	// Initialize GORM
	db := database.InitDB()

	// Get underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		logger.Log(logger.ERROR, "Failed to get underlying *sql.DB: "+err.Error())
	}

	// Close connection when main function ends
	defer sqlDB.Close()

	initServer(db)
}

func initServer(db *gorm.DB) {
	// Auto migrate database
	logger.Log(logger.DEBUG, "Running database migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Meal{},
		&models.UserStats{},
	); err != nil {
		logger.Log(logger.ERROR, "Failed to migrate database: "+err.Error())
		return
	}
	// create router with gorm db
	router := router.NewRouter(db)
	ctx := context.Background()
	// enable cors
	router.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		}
	}())
	if err := seed.SeedDatabase(db, ctx); err != nil {
		logger.Log(logger.ERROR, "Error seeding database :: "+err.Error())
	}

	// router middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	logger.Log(logger.DEBUG, "Starting server on port 8080")

	// register custom validator in gin validation engine
	// context is intrinsic from ctx
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("boolean", validators.ValidateBoolean)
	}
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	srv.ListenAndServe()
}
