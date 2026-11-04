package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	authHTTPConfig, err := config.LoadAuthHTTPConfig()
	if err != nil {
		log.Fatalf("Failed to load auth HTTP config: %v", err)
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
	authHandler := handler.NewAuthHandler(authService, authHTTPConfig)

	healthHandler := handler.NewHealthHandler(db)

	r := gin.Default()
	if err := r.SetTrustedProxies(authHTTPConfig.TrustedProxies); err != nil {
		log.Fatalf("Failed to configure trusted proxies: %v", err)
	}
	r.Use(middleware.BrowserSecurity(authHTTPConfig.AllowedOrigins))
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", healthHandler.CheckHealth)

		authHandler.RegisterRoutes(apiV1.Group("/auth"), authConfig.JWTSecret)

		apiV1.GET("/plants", plantHandler.GetPlants)
		apiV1.GET("/diseases", diseaseHandler.GetDiseases)
		apiV1.GET("/diseases/:class_name", diseaseHandler.GetDiseaseByClassName)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Server on port %s", port)
	server := &http.Server{
		Addr: ":" + port, Handler: r,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second,
		WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
