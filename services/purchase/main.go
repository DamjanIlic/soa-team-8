package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"purchase/handler"
	"purchase/middleware"
	"purchase/model"
	"purchase/repo"
	"purchase/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	host := getEnv("DB_HOST", "purchase-db")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "super")
	dbname := getEnv("DB_NAME", "purchasedb")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Belgrade",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// migracije
	db.AutoMigrate(&model.ShoppingCart{})
	db.AutoMigrate(&model.OrderItem{})
	db.AutoMigrate(&model.TourPurchaseToken{})

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

	// repo
	cartRepo := &repo.CartRepository{DB: db}
	itemRepo := &repo.ItemRepository{DB: db}
	tokenRepo := &repo.TokenRepository{DB: db}

	// services
	cartService := &service.CartService{
		CartRepo: cartRepo,
		ItemRepo: itemRepo,
	}
	tokenService := &service.TokenService{
		CartRepo:  cartRepo,
		TokenRepo: tokenRepo,
		ItemRepo:  itemRepo,
	}

	// handlers
	cartHandler := &handler.CartHandler{CartService: cartService}
	tokenHandler := &handler.TokenHandler{TokenService: tokenService}

	startServer(cartHandler, tokenHandler)
}

func startServer(cartHandler *handler.CartHandler, tokenHandler *handler.TokenHandler) {
	router := mux.NewRouter().StrictSlash(true)

	// JWT middleware
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Purchase service is alive 🚀"))
	}).Methods("GET")

	// Cart endpoints
	api.HandleFunc("/cart", cartHandler.CreateCart).Methods("POST")
	api.HandleFunc("/cart", cartHandler.GetCart).Methods("GET")
	api.HandleFunc("/cart/items", cartHandler.AddItem).Methods("POST")
	api.HandleFunc("/cart/items/{itemId}", cartHandler.RemoveItem).Methods("DELETE")
	api.HandleFunc("/cart/total", cartHandler.GetTotal).Methods("GET")

	// Token endpoints (checkout)
	api.HandleFunc("/cart/checkout", tokenHandler.Checkout).Methods("POST")

	port := getEnv("PORT", "8080")
	log.Printf("Purchase service starting on :%s 🚀\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
