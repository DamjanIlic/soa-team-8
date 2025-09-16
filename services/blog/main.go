package main

import (
	"blog/handler"
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
	api.HandleFunc("/blogs", blogHandler.Create).Methods("POST")
	api.HandleFunc("/blogs", blogHandler.GetAll).Methods("GET")
	api.HandleFunc("/blogs/{id}", blogHandler.Get).Methods("GET")
	api.HandleFunc("/blogs/{id}/like", blogHandler.Like).Methods("POST")
	api.HandleFunc("/blogs/{id}/unlike", blogHandler.Unlike).Methods("POST")

	// Comment routes
	api.HandleFunc("/blogs/{id}/comments", commentHandler.Create).Methods("POST")
	api.HandleFunc("/blogs/{id}/comments", commentHandler.GetByBlogID).Methods("GET")

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Blog service running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
