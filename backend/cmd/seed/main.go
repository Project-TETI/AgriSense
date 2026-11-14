package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Project-TETI/AgriSense/backend/internal/catalogseed"
	"github.com/Project-TETI/AgriSense/backend/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	dir := flag.String("dir", "data/catalog", "directory containing catalog JSON files")
	flag.Parse()
	plants, err := catalogseed.Load(*dir)
	if err != nil {
		log.Fatal(err)
	}
	_ = godotenv.Load()
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("database initialization failed")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("database connection unavailable")
	}
	defer sqlDB.Close()
	counts, err := catalogseed.Apply(db, plants)
	if err != nil {
		log.Fatal("catalog seed failed; transaction rolled back")
	}
	fmt.Printf("Plants: %d created, %d skipped; diseases: %d created, %d skipped\n", counts.PlantsCreated, counts.PlantsSkipped, counts.DiseasesCreated, counts.DiseasesSkipped)
}
