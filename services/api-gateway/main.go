package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	handler := enableCORS(r)

	// Proxy ka TourService
	r.PathPrefix("/api/tours").Handler(proxy("http://tour-service:8080"))

	// Proxy ka AuthService
	r.PathPrefix("/api/auth").Handler(proxy("http://auth-service:8080"))

	// Proxy ka BlogService
	r.PathPrefix("/api/blogs").Handler(proxy("http://blog-service:8080"))

	// Proxy ka FollowService
	r.PathPrefix("/api/follow").Handler(proxy("http://follow-service:8000"))

	// Proxy ka StakeholdersService
	r.PathPrefix("/api/stakeholders").Handler(proxy("http://stakeholders-service:8080"))

	log.Println("API Gateway is running on :8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal("Gateway failed: ", err)
	}
}

// proxy pravi reverse proxy ka mikroservisu
func proxy(target string) http.Handler {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Cannot parse URL %s: %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(url)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONS preflight zahtev
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
