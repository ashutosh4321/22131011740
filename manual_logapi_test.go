package main

import (
	"os"
	"log"
	"encoding/json"
	"bytes"
	"net/http"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	// Set up test data
	stack := "backend"
	level := "info"
	pkg := "service"
	message := "Testing middleware"

	api_url := os.Getenv("LOG_API_URL")
	api_token := os.Getenv("LOG_API_ACCESS_TOKEN")
	if api_url == "" || api_token == "" {
		log.Fatal("log api url or token is not found in env")
	}

	jsonData := map[string]string{
		"stack": stack,
		"level": level,
		"pkg": pkg,
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
	} else {
		log.Println("Log API call succeeded with status 200")
	}
} 