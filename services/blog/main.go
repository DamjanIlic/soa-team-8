package main

import (
	"blog/handler"
	"blog/model"
	"blog/repo"
	"blog/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	if dbName == "" {
		dbName = "blog"
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("DB connection failed, retrying in 3s... (%d/10)\n", i+1)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	// Auto-migrate tabele
	err = db.AutoMigrate(&model.Blog{}, &model.Comment{})
	if err != nil {
		log.Fatal("Failed to auto-migrate tables: ", err)
	}

	blogRepo := &repo.BlogRepository{DatabaseConnection: db}
	blogService := &service.BlogService{BlogRepo: blogRepo}
	blogHandler := &handler.BlogHandler{BlogService: blogService}

	commentRepo := &repo.CommentRepository{DB: db}
	commentService := &service.CommentService{CommentRepo: commentRepo}
	commentHandler := &handler.CommentHandler{CommentService: commentService}

	r := mux.NewRouter()
	r.HandleFunc("/blogs", blogHandler.Create).Methods("POST")
	r.HandleFunc("/blogs/{id}", blogHandler.Get).Methods("GET")
	r.HandleFunc("/blogs/{id}/like", blogHandler.Like).Methods("POST")
	r.HandleFunc("/blogs/{id}/unlike", blogHandler.Unlike).Methods("POST")

	//comment
	r.HandleFunc("/blogs/{id}/comments", commentHandler.Create).Methods("POST")
	r.HandleFunc("/blogs/{id}/comments", commentHandler.GetByBlogID).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Blog service running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
