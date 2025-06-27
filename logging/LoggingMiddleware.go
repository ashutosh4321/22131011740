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
	if api_url == "" {
		log.Fatal("log api url is not found in env")
	}

	//Get new access token
	authUrl := os.Getenv("LOG_AUTH_URL")
	authPayload := map[string]string{
		"email":      os.Getenv("LOG_EMAIL"),
		"name":       os.Getenv("LOG_NAME"),
		"rollNo":     os.Getenv("LOG_ROLLNO"),
		"accessCode": os.Getenv("LOG_ACCESS_CODE"),
		"clientID":   os.Getenv("LOG_CLIENT_ID"),
		"clientSecret": os.Getenv("LOG_CLIENT_SECRET"),
	}
	authBytes, err := json.Marshal(authPayload)
	if err != nil {
		log.Printf("Failed to marshal auth payload: %v", err)
		return
	}

	// TODO: Implement token caching

	authReq, err := http.NewRequest("POST", authUrl, bytes.NewBuffer(authBytes))
	if err != nil {
		log.Printf("Failed to create auth request: %v", err)
		return
	}
	authReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	authResp, err := client.Do(authReq)
	if err != nil {
		log.Printf("Failed to send auth request: %v", err)
		return
	}
	defer authResp.Body.Close()
	if authResp.StatusCode != http.StatusOK {
		log.Printf("Auth API returned non-200 status: %d", authResp.StatusCode)
		return
	}
	var authRespData struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(authResp.Body).Decode(&authRespData); err != nil {
		log.Printf("Failed to decode auth response: %v", err)
		return
	}
	api_token := authRespData.AccessToken

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