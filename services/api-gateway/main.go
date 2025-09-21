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

	// Proxy ka PurchaseService
	r.PathPrefix("/api/cart").Handler(proxy("http://purchase-service:8080"))

	log.Println("API Gateway is running on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
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
