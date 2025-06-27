package logging

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"bytes"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func LogApi(stack string, level string, pkg string, message string){
	godotenv.Load(".env")

	api_url := os.Getenv("LOG_API_URL")
	api_token := os.Getenv("LOG_API_ACCESS_TOKEN")
	if api_url == "" || api_token == "" {
		log.Fatal("log api url or token is not found in env")
	}

	jsonData := map[string]string{
		"stack": stack,
		"level": level,
		"package": pkg,
		"message": message,
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Printf("Failed to marshal log data: %v", err)
		return
	}

	req, err := http.NewRequest("POST", api_url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("Failed to create log API request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send log API request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Log API returned non-200 status: %d", resp.StatusCode)
	}
}

func CustomLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		LogApi("backend", "info", "service", "Received request: " + r.Method + " " + r.URL.Path)
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Wrap the ResponseWriter to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		LogApi("backend", "info", "service", "Sent response: " + r.Method + " " + r.URL.Path + " - Status: " + strconv.Itoa(ww.Status()))
		log.Printf("Sent response: %s %s - Status: %d", r.Method, r.URL.Path, ww.Status())
	})
}