package main

import (
	"auth/handler"
	"auth/model"
	"auth/proto"
	"auth/repo"
	"auth/service"
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// gRPC server samo za Login i GetUser
type UserServer struct {
	proto.UnimplementedAuthServiceServer
	UserService *service.UserService
}

func (s *UserServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.UserService.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{Token: token}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := s.UserService.UserRepo.FindByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserResponse{
		Id:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		Blocked:  user.Blocked,
	}, nil
}

func main() {
	db := initDB()

	userRepo := &repo.UserRepo{DB: db}
	userService := &service.UserService{UserRepo: userRepo}
	userHandler := &handler.UserHandler{UserService: userService}

	// --- REST rute (ostaju iste) ---
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/auth/block/{id}", userHandler.BlockUser).Methods("POST")

	// HTTP server za REST
	httpPort := getEnv("PORT", "8086")
	go func() {
		log.Printf("HTTP REST server running on :%s\n", httpPort)
		log.Fatal(http.ListenAndServe(":"+httpPort, router))
	}()

	// --- gRPC server za Login i GetUser ---
	grpcPort := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", "0.0.0.0:"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, &UserServer{UserService: userService})

	log.Printf("gRPC server running on :%s\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}

// --- DB konekcija ---
func initDB() *gorm.DB {
	host := getEnv("DB_HOST", "stakeholders-db")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "super")
	dbname := getEnv("DB_NAME", "stakeholdersdb")
	port := getEnv("DB_PORT", "5432")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
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
