package main

import (
	"log"
	"os"

	"github.com/Project-TETI/AgriSense/backend/internal/config"
	"github.com/Project-TETI/AgriSense/backend/internal/handler"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	db := config.InitDB()

	plantRepo := repository.NewPlantRepository(db)
	plantService := service.NewPlantService(plantRepo)
	plantHandler := handler.NewPlantHandler(plantService)

	diseaseRepo := repository.NewDiseaseRepository(db)
	diseaseService := service.NewDiseaseService(diseaseRepo)
	diseaseHandler := handler.NewDiseaseHandler(diseaseService)

	healthHandler := handler.NewHealthHandler(db)

	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", healthHandler.CheckHealth)
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
