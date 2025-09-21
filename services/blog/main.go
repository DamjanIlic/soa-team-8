package main

import (
	"blog/handler"
	"blog/middleware"
	"blog/repo"
	"blog/service"
	"log"
	"net/http"
	"os"

	"context"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongo() *mongo.Database {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "blog")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	return client.Database(dbName)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	db := initMongo()

	// Repositories
	blogRepo := &repo.BlogRepository{Collection: db.Collection("blogs")}
	commentRepo := &repo.CommentRepository{Collection: db.Collection("comments")}
	likeRepo := &repo.LikeRepository{Collection: db.Collection("likes")}

	// Services
	blogService := &service.BlogService{BlogRepo: blogRepo, LikeRepo: likeRepo}
	commentService := &service.CommentService{CommentRepo: commentRepo}

	// Handlers
	blogHandler := &handler.BlogHandler{BlogService: blogService}
	commentHandler := &handler.CommentHandler{CommentService: commentService}

	// Router
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	// Blog routes
	api.HandleFunc("/blogs", blogHandler.GetAll).Methods("GET")
	api.HandleFunc("/blogs/{id}", blogHandler.Get).Methods("GET")

	// Kreiranje bloga i like/unlike zahtevaju autentifikaciju
	api.Handle("/blogs", middleware.JWTMiddleware(http.HandlerFunc(blogHandler.Create))).Methods("POST")
	api.Handle("/blogs/{id}/like", middleware.JWTMiddleware(http.HandlerFunc(blogHandler.Like))).Methods("POST")
	api.Handle("/blogs/{id}/unlike", middleware.JWTMiddleware(http.HandlerFunc(blogHandler.Unlike))).Methods("POST")

	// Comment routes (ako želiš, kreiranje komentara može takođe zahtevati JWT)
	api.HandleFunc("/blogs/{id}/comments", commentHandler.GetByBlogID).Methods("GET")
	api.Handle("/blogs/{id}/comments", middleware.JWTMiddleware(http.HandlerFunc(commentHandler.Create))).Methods("POST")

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Blog service running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
