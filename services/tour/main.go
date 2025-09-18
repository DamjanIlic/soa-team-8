package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tour/handler"
	"tour/middleware"
	"tour/model"
	"tour/repo"
	"tour/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	host := getEnv("DB_HOST", "tour-db")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "super")
	dbname := getEnv("DB_NAME", "tourdb")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Belgrade",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migracije
	db.AutoMigrate(&model.Tour{})
	db.AutoMigrate(&model.KeyPoint{})
	db.AutoMigrate(&model.Duration{})

	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	db := initDB()

	// Repozitorijumi
	tourRepo := &repo.TourRepository{DatabaseConnection: db}
	keyPointRepo := &repo.KeyPointRepository{DatabaseConnection: db}
	durationRepo := &repo.DurationRepository{DatabaseConnection: db}

	// Servisi
	tourService := &service.TourService{TourRepo: tourRepo}
	keyPointService := &service.KeyPointService{
		KeyPointRepo: keyPointRepo,
		TourRepo:     tourRepo,
	}
	durationService := &service.DurationService{
		DurationRepo: durationRepo,
		TourRepo:     tourRepo,
	}

	// Handleri
	tourHandler := &handler.TourHandler{TourService: tourService}
	keyPointHandler := &handler.KeyPointHandler{KeyPointService: keyPointService}
	durationHandler := &handler.DurationHandler{DurationService: durationService}

	startServer(tourHandler, keyPointHandler, durationHandler)
}

func startServer(tourHandler *handler.TourHandler, keyPointHandler *handler.KeyPointHandler, durationHandler *handler.DurationHandler) {
	router := mux.NewRouter().StrictSlash(true)

	// JWT middleware na svim API rutama
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	// Tour endpoints
	api.HandleFunc("/tours", tourHandler.CreateTour).Methods("POST")
	api.HandleFunc("/tours/{id}", tourHandler.GetTour).Methods("GET")
	api.HandleFunc("/tours", tourHandler.GetAllTours).Methods("GET")
	api.HandleFunc("/tours/authors/{authorId}", tourHandler.GetToursByAuthor).Methods("GET")
	api.HandleFunc("/tours/{id}/publish", tourHandler.PublishTour).Methods("POST")
	api.HandleFunc("/tours/{id}/archive", tourHandler.ArchiveTour).Methods("POST")
	api.HandleFunc("/tours/{id}/reactivate", tourHandler.ReactivateTour).Methods("POST")

	// KeyPoint endpoints
	api.HandleFunc("/tours/{tourId}/keypoints", keyPointHandler.CreateKeyPoint).Methods("POST")
	api.HandleFunc("/tours/{tourId}/keypoints", keyPointHandler.GetKeyPointsByTour).Methods("GET")
	api.HandleFunc("/tours/keypoints/{id}", keyPointHandler.GetKeyPoint).Methods("GET")
	api.HandleFunc("/tours/keypoints/{id}", keyPointHandler.UpdateKeyPoint).Methods("PUT")
	api.HandleFunc("/tours/keypoints/{id}", keyPointHandler.DeleteKeyPoint).Methods("DELETE")

	// Duration endpoints
	api.HandleFunc("/tours/{tourId}/durations", durationHandler.AddDuration).Methods("POST")
	api.HandleFunc("/tours/{tourId}/durations", durationHandler.GetDurationsByTour).Methods("GET")

	// Staticki fajlovi
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	port := getEnv("PORT", "8080")
	log.Printf("Tour service starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
