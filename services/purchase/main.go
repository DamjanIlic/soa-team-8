package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"purchase/grpc"
	"purchase/handler"
	"purchase/middleware"
	"purchase/model"
	"purchase/pb"
	"purchase/repo"
	"purchase/service"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcServer "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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

	// handlers (samo za ne-RPC endpoints)
	cartHandler := &handler.CartHandler{CartService: cartService}
	tokenHandler := &handler.TokenHandler{TokenService: tokenService}

	// gRPC server
	grpcPurchaseServer := &grpc.PurchaseGRPCServer{
		TokenService: tokenService,
		CartService:  cartService,
	}

	// Pokretanje gRPC servera u goroutine
	go startGRPCServer(grpcPurchaseServer)

	// Pokretanje HTTP servera sa gRPC-Gateway
	startHTTPServerWithGateway(cartHandler, tokenHandler)
}

func startGRPCServer(purchaseServer *grpc.PurchaseGRPCServer) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	s := grpcServer.NewServer()
	pb.RegisterPurchaseServiceServer(s, purchaseServer)

	reflection.Register(s)

	log.Println("gRPC Purchase service running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func startHTTPServerWithGateway(cartHandler *handler.CartHandler, tokenHandler *handler.TokenHandler) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// gRPC-Gateway mux
	gwMux := runtime.NewServeMux()
	opts := []grpcServer.DialOption{grpcServer.WithTransportCredentials(insecure.NewCredentials())}

	// Registruj gRPC-Gateway
	if err := pb.RegisterPurchaseServiceHandlerFromEndpoint(ctx, gwMux, "localhost:50051", opts); err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Gorilla mux za ostale endpoints
	router := mux.NewRouter().StrictSlash(true)

	// JWT middleware za sve API rute
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	// Health check endpoint
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Purchase service is alive 🚀"))
	}).Methods("GET")

	// RPC endpoints se sada serviraju preko gRPC-Gateway
	// /api/cart (GetCart) - automatski preko gRPC
	// /api/cart/checkout (Checkout) - automatski preko gRPC

	// Ostali HTTP endpoints (koji nisu RPC)
	api.HandleFunc("/cart", cartHandler.CreateCart).Methods("POST")
	api.HandleFunc("/cart/items", cartHandler.AddItem).Methods("POST")
	api.HandleFunc("/cart/items/{itemId}", cartHandler.RemoveItem).Methods("DELETE")
	api.HandleFunc("/cart/total", cartHandler.GetTotal).Methods("GET")

	// Token endpoints
	api.HandleFunc("/cart/tokens/purchased", tokenHandler.GetPurchasedTours).Methods("GET")
	api.HandleFunc("/cart/tokens/{tokenId}/executed", tokenHandler.MarkAsExecuted).Methods("PUT")
	api.HandleFunc("/cart/tokens/{tokenId}/reviewed", tokenHandler.MarkAsReviewed).Methods("PUT")
	api.HandleFunc("/cart/tokens/tour/{tourId}", tokenHandler.GetTokenByTouristAndTour).Methods("GET")

	// Kombinuj Gorilla router sa gRPC-Gateway
	router.PathPrefix("/").Handler(gwMux)

	port := getEnv("PORT", "8080")
	log.Printf("HTTP Purchase service with gRPC-Gateway starting on :%s 🚀\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
