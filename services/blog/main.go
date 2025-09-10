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

func initDB() *gorm.DB {
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "blog")
	dbPort := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Belgrade",
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

	if err := db.AutoMigrate(&model.Blog{}); err != nil {
		log.Fatal("Failed to auto-migrate Blog table: ", err)
	}

	return db
}

// func main() {
// 	database := initDB()
// 	db := database
// 	// blogRepo := &repo.BlogRepository{DatabaseConnection: database}
// 	// // Auto-migrate tabele
// 	// err = db.AutoMigrate(&model.Blog{}, &model.Comment{})
// 	// if err != nil {
// 	// 	log.Fatal("Failed to auto-migrate tables: ", err)
// 	// }

// 	if err := db.AutoMigrate(&model.Blog{}); err != nil {
// 		log.Fatal("Failed to auto-migrate Blog table: ", err)
// 	}

// 	return db
// }

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	database := initDB()
	db := database
	blogRepo := &repo.BlogRepository{DatabaseConnection: database}
	blogService := &service.BlogService{BlogRepo: blogRepo}
	blogHandler := &handler.BlogHandler{BlogService: blogService}

	commentRepo := &repo.CommentRepository{DB: db}
	commentService := &service.CommentService{CommentRepo: commentRepo}
	commentHandler := &handler.CommentHandler{CommentService: commentService}

	r := mux.NewRouter()
	r.HandleFunc("/blogs", blogHandler.Create).Methods("POST")
	r.HandleFunc("/blogs", blogHandler.GetAll).Methods("GET")
	r.HandleFunc("/blogs/{id}", blogHandler.Get).Methods("GET")
	r.HandleFunc("/blogs/{id}/like", blogHandler.Like).Methods("POST")
	r.HandleFunc("/blogs/{id}/unlike", blogHandler.Unlike).Methods("POST")

	//port := getEnv("PORT", "8080")
	//comment
	r.HandleFunc("/blogs/{id}/comments", commentHandler.Create).Methods("POST")
	r.HandleFunc("/blogs/{id}/comments", commentHandler.GetByBlogID).Methods("GET")

	port := getEnv("PORT", "8080")

	log.Printf("Blog service running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
