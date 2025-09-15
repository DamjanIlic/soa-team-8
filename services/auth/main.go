package main

import (
	"auth-service/handler"
	"auth-service/model"
	"auth-service/repo"
	"auth-service/service"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db := initDB()

	userRepo := &repo.UserRepo{DB: db}
	userService := &service.UserService{UserRepo: userRepo}
	userHandler := &handler.UserHandler{UserService: userService}

	// Router sa /api prefiksom
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	// Rute za auth-service sa prefiksom /api
	api.HandleFunc("/auth/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", userHandler.Login).Methods("POST")
	api.HandleFunc("/auth/block/{id}", userHandler.BlockUser).Methods("POST")

	port := getEnv("PORT", "8080")
	log.Printf("Auth service running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initDB() *gorm.DB {
	host := getEnv("DB_HOST", "auth-db")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "super")
	dbname := getEnv("DB_NAME", "authdb")
	port := getEnv("DB_PORT", "5432")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.User{})
	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
