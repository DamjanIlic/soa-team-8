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

	db.AutoMigrate(&model.Tour{})
	db.AutoMigrate(&model.KeyPoint{})
	db.AutoMigrate(&model.Review{})

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

	tourRepo := &repo.TourRepository{DatabaseConnection: db}
	keyPointRepo := &repo.KeyPointRepository{DatabaseConnection: db}
	reviewRepo := &repo.ReviewRepository{DatabaseConnection: db}

	tourService := &service.TourService{TourRepo: tourRepo}
	keyPointService := &service.KeyPointService{
		KeyPointRepo: keyPointRepo,
		TourRepo:     tourRepo,
	}
	reviewService := &service.ReviewService{
		ReviewRepo: reviewRepo,
	}

	tourHandler := &handler.TourHandler{TourService: tourService}
	keyPointHandler := &handler.KeyPointHandler{KeyPointService: keyPointService}
	reviewHandler := &handler.ReviewHandler{ReviewService: reviewService}

	startServer(tourHandler, keyPointHandler, reviewHandler)
}

func startServer(tourHandler *handler.TourHandler, keyPointHandler *handler.KeyPointHandler, reviewHandler *handler.ReviewHandler) {
	router := mux.NewRouter().StrictSlash(true)

	// JWT middleware na svim API rutama
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	// Tour endpoints
	api.HandleFunc("/tours", tourHandler.CreateTour).Methods("POST")
	api.HandleFunc("/tours/{id}", tourHandler.GetTour).Methods("GET")
	api.HandleFunc("/tours", tourHandler.GetAllTours).Methods("GET")
	api.HandleFunc("/tours/authors/{authorId}", tourHandler.GetToursByAuthor).Methods("GET")

	// KeyPoint endpoints
	api.HandleFunc("/tours/{tourId}/keypoints", keyPointHandler.CreateKeyPoint).Methods("POST")
	api.HandleFunc("/tours/{tourId}/keypoints", keyPointHandler.GetKeyPointsByTour).Methods("GET")
	api.HandleFunc("/tours/keypoints/{id}", keyPointHandler.GetKeyPoint).Methods("GET")
	api.HandleFunc("/tours/keypoints/{id}", keyPointHandler.UpdateKeyPoint).Methods("PUT")
	api.HandleFunc("/tours/keypoints/{id}", keyPointHandler.DeleteKeyPoint).Methods("DELETE")

	// Review endpoints
	api.HandleFunc("/tours/{tourId}/reviews", reviewHandler.CreateReview).Methods("POST")
	api.HandleFunc("/tours/{tourId}/reviews", reviewHandler.GetReviewsByTour).Methods("GET")

	// Staticki fajlovi
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	port := getEnv("PORT", "8080")
	log.Printf("Tour service starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
