package main

import (
	"log"
	"net/http"
	"stakeholder/handler"
	"stakeholder/model"
	"stakeholder/repo"
	"stakeholder/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=super dbname=stakeholdersdb port=5432 sslmode=disable TimeZone=Europe/Belgrade"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate tabele
	db.AutoMigrate(&model.Stakeholder{})
	db.AutoMigrate(&model.User{})

	return db
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

	//admin
	api.HandleFunc("/admin/users", userHandler.GetAllUsers).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
