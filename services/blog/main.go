package main

import (
	"blog/handler"
	"blog/metrics"
	"blog/middleware"
	"blog/repo"
	"blog/service"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer(serviceName string) (func(context.Context) error, error) {
	// Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger:14268/api/traces")))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp,
			sdktrace.WithMaxExportBatchSize(1),       // pošalji span odmah
			sdktrace.WithBatchTimeout(1*time.Second), // ili 1s timeout
		),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("blog-service"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

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
	shutdown, err := initTracer("blog-service")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			log.Printf("error shutting down tracer: %v", err)
		}
	}()
	tracer := otel.Tracer("blog-service")
	metrics.Init()
	// Repositories
	blogRepo := &repo.BlogRepository{Collection: db.Collection("blogs")}
	commentRepo := &repo.CommentRepository{Collection: db.Collection("comments")}
	likeRepo := &repo.LikeRepository{Collection: db.Collection("likes")}

	// Services
	stakeholderClient := service.NewStakeholderClient("http://stakeholders-service:8080")
	blogService := &service.BlogService{BlogRepo: blogRepo, LikeRepo: likeRepo}
	commentService := &service.CommentService{
		CommentRepo:       commentRepo,
		StakeholderClient: stakeholderClient,
	}

	// Handlers
	blogHandler := &handler.BlogHandler{BlogService: blogService}
	commentHandler := &handler.CommentHandler{CommentService: commentService}

	// Router
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	// Blog routes
	//api.HandleFunc("/blogs", blogHandler.GetAll).Methods("GET")
	api.HandleFunc("/blogs", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "GetAllBlogs")
		defer span.End()
		blogHandler.GetAll(w, r.WithContext(ctx))
	}).Methods("GET")
	api.HandleFunc("/blogs/{id}", blogHandler.Get).Methods("GET")

	// Kreiranje bloga i like/unlike zahtevaju autentifikaciju
	api.Handle("/blogs", middleware.JWTMiddleware(http.HandlerFunc(blogHandler.Create))).Methods("POST")
	api.Handle("/blogs/{id}/like", middleware.JWTMiddleware(http.HandlerFunc(blogHandler.Like))).Methods("POST")
	api.Handle("/blogs/{id}/unlike", middleware.JWTMiddleware(http.HandlerFunc(blogHandler.Unlike))).Methods("POST")

	// Comment routes (ako želiš, kreiranje komentara može takođe zahtevati JWT)
	api.HandleFunc("/blogs/{id}/comments", commentHandler.GetByBlogID).Methods("GET")
	api.Handle("/blogs/{id}/comments", middleware.JWTMiddleware(http.HandlerFunc(commentHandler.Create))).Methods("POST")

	//prometheus
	router.Handle("/metrics", promhttp.Handler())

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Blog service running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
