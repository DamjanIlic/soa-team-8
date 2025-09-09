package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"stakeholder/handler"
	"stakeholder/model"
	"stakeholder/repo"
	"stakeholder/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	// čitaj konfiguraciju iz environmenta, sa fallback vrednostima
	host := getEnv("DB_HOST", "stakeholders-db")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "super")
	dbname := getEnv("DB_NAME", "stakeholdersdb")
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
	db.AutoMigrate(&model.Stakeholder{})
	db.AutoMigrate(&model.User{})

	return db
}

// getEnv čita env var, fallback ako nije definisana
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	database := initDB()

	stakeholderRepo := &repo.StakeholderRepository{DatabaseConnection: database}
	stakeholderService := &service.StakeholderService{StakeholderRepo: stakeholderRepo}
	stakeholderHandler := &handler.StakeholderHandler{StakeholderService: stakeholderService}

	userRepo := &repo.UserRepo{DatabaseConnection: database}
	userService := &service.UserService{UserRepo: userRepo}
	userHandler := &handler.UserHandler{UserService: userService}

	startServer(stakeholderHandler, userHandler)
}

func startServer(handler *handler.StakeholderHandler, userHandler *handler.UserHandler) {
	router := mux.NewRouter().StrictSlash(true)

	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/stakeholders/{id}", handler.Get).Methods("GET")
	api.HandleFunc("/stakeholders", handler.Create).Methods("POST")

	// admin endpoint
	api.HandleFunc("/admin/users", userHandler.GetAllUsers).Methods("GET")

	// endpoint za registraciju neregistrovanih korisnika
	api.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")

	// static fajlovi
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	// port iz env varijable, fallback 8080
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
