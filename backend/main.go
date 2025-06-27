package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/ashutosh4321/22131011740/logging"
)

// Request and response structs for short URL creation
type CreateShortURLRequest struct {
	URL       string `json:"url"`
	Validity  int    `json:"validity,omitempty"`
	Shortcode string `json:"shortcode,omitempty"`
}

type CreateShortURLResponse struct {
	ShortLink string `json:"shortLink"`
	Expiry    string `json:"expiry"`
}

// Response struct for short URL statistics
type ShortURLStatsResponse struct {
	Shortcode   string        `json:"shortcode"`
	OriginalURL string        `json:"originalUrl"`
	CreatedAt   string        `json:"createdAt"`
	Expiry      string        `json:"expiry"`
	Clicks      int           `json:"clicks"`
	ClickData   []ClickDetail `json:"clickData"`
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in env")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	router.Use(logging.CustomLogger)

	router.Post("/shorturls", createShortURLHandler)
	router.Get("/shorturls/{shortcode}", getShortURLStatsHandler)
	router.Get("/{shortcode}", redirectHandler)

	srv := &http.Server{
		Handler: router,
		Addr: ":"+port,
	}

	log.Printf("Server Staring on Port %v", port)
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", port)
}
