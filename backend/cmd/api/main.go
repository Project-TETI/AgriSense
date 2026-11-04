package main

import (
	"log"
	"os"

	"github.com/Project-TETI/AgriSense/backend/internal/config"
	"github.com/Project-TETI/AgriSense/backend/internal/handler"
	"github.com/Project-TETI/AgriSense/backend/internal/middleware"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		log.Fatalf("Failed to load auth config: %v", err)
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	plantRepo := repository.NewPlantRepository(db)
	plantService := service.NewPlantService(plantRepo)
	plantHandler := handler.NewPlantHandler(plantService)

	diseaseRepo := repository.NewDiseaseRepository(db)
	diseaseService := service.NewDiseaseService(diseaseRepo)
	diseaseHandler := handler.NewDiseaseHandler(diseaseService)

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authTransaction := repository.NewAuthTransaction(db)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, authTransaction, authConfig.JWTSecret, authConfig.AccessDuration, authConfig.RefreshDuration)
	authHandler := handler.NewAuthHandler(authService)

	healthHandler := handler.NewHealthHandler(db)

	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", healthHandler.CheckHealth)

		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)

			protectedAuth := authGroup.Group("")
			protectedAuth.Use(middleware.AuthMiddleware(authConfig.JWTSecret))
			{
				protectedAuth.POST("/logout", authHandler.Logout)
				protectedAuth.GET("/me", authHandler.GetMe)
			}
		}

		apiV1.GET("/plants", plantHandler.GetPlants)
		apiV1.GET("/diseases", diseaseHandler.GetDiseases)
		apiV1.GET("/diseases/:class_name", diseaseHandler.GetDiseaseByClassName)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
