package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tour/handler"
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

	// Auto migrate tabele
	db.AutoMigrate(&model.Tour{})
	db.AutoMigrate(&model.KeyPoint{})

	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	database := initDB()

	// Repository instances
	tourRepo := &repo.TourRepository{DatabaseConnection: database}
	keyPointRepo := &repo.KeyPointRepository{DatabaseConnection: database}

	// Service instances
	tourService := &service.TourService{TourRepo: tourRepo}
	keyPointService := &service.KeyPointService{
		KeyPointRepo: keyPointRepo,
		TourRepo:     tourRepo,
	}

	// Handler instances
	tourHandler := &handler.TourHandler{TourService: tourService}
	keyPointHandler := &handler.KeyPointHandler{KeyPointService: keyPointService}

	startServer(tourHandler, keyPointHandler)
}

func startServer(tourHandler *handler.TourHandler, keyPointHandler *handler.KeyPointHandler) {
	router := mux.NewRouter().StrictSlash(true)

	api := router.PathPrefix("/api").Subrouter()

	// Tour endpoints
	api.HandleFunc("/tours", tourHandler.CreateTour).Methods("POST")
	api.HandleFunc("/tours/{id}", tourHandler.GetTour).Methods("GET")
	api.HandleFunc("/tours", tourHandler.GetAllTours).Methods("GET")
	api.HandleFunc("/authors/{authorId}/tours", tourHandler.GetToursByAuthor).Methods("GET")

	// KeyPoint endpoints
	api.HandleFunc("/tours/{tourId}/keypoints", keyPointHandler.CreateKeyPoint).Methods("POST")
	api.HandleFunc("/tours/{tourId}/keypoints", keyPointHandler.GetKeyPointsByTour).Methods("GET")
	api.HandleFunc("/keypoints/{id}", keyPointHandler.GetKeyPoint).Methods("GET")
	api.HandleFunc("/keypoints/{id}", keyPointHandler.UpdateKeyPoint).Methods("PUT")
	api.HandleFunc("/keypoints/{id}", keyPointHandler.DeleteKeyPoint).Methods("DELETE")

	// static fajlovi
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	port := getEnv("PORT", "8080")
	log.Printf("Tour service starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}