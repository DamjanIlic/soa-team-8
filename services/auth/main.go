package main

import (
	"auth/model"
	"auth/repo"
	"auth/service"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	gw "auth/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// UserServer implementira AuthServiceServer
type UserServer struct {
	gw.UnimplementedAuthServiceServer
	UserService *service.UserService
}

// Login RPC metoda
func (s *UserServer) Login(ctx context.Context, req *gw.LoginRequest) (*gw.AuthResponse, error) {
	token, err := s.UserService.Login(req.Email, req.Password)
	if err != nil {
		return &gw.AuthResponse{Message: err.Error()}, nil
	}
	return &gw.AuthResponse{
		Token:   token,
		Message: "login successful",
	}, nil
}

// GetUser RPC metoda
func (s *UserServer) GetUser(ctx context.Context, req *gw.GetUserRequest) (*gw.GetUserResponse, error) {
	user, err := s.UserService.GetUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &gw.GetUserResponse{
		Id:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		Blocked:  user.Blocked,
	}, nil
}

// HTTP handler za registraciju (bez RPC)
func RegisterHandler(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		user := &model.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
			Role:     model.Role(req.Role),
		}

		if err := svc.RegisterUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "registration successful",
		})
	}
}

func main() {
	db := initDB()
	userRepo := &repo.UserRepo{DB: db}
	userService := &service.UserService{UserRepo: userRepo}

	grpcPort := getEnv("GRPC_PORT", "50051")
	httpPort := getEnv("PORT", "8080")

	// gRPC server
	grpcServer := grpc.NewServer()
	gw.RegisterAuthServiceServer(grpcServer, &UserServer{UserService: userService})

	// Omogući reflection za grpcurl
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("gRPC Auth service running on :%s", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// gRPC-Gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := gw.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "auth-service:"+grpcPort, opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// HTTP mux kombinuje Register handler i ostale RPC endpoint-e
	httpMux := http.NewServeMux()
	httpMux.Handle("/api/auth/register", RegisterHandler(userService))
	httpMux.Handle("/", mux)

	log.Printf("HTTP Auth service running on :%s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, httpMux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

func initDB() *gorm.DB {
	host := getEnv("DB_HOST", "stakeholders-db")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "super")
	dbname := getEnv("DB_NAME", "stakeholdersdb")
	port := getEnv("DB_PORT", "5432")

	dsn := "host=" + host + " user=" + user + " password=" + password +
		" dbname=" + dbname + " port=" + port + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
