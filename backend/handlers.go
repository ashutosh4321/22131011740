package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"strings"
	"github.com/go-chi/chi/v5"
	"github.com/ashutosh4321/22131011740/logging"
)

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateShortURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.LogApi("backend", "error", "handler", "Invalid request body: "+err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	logging.LogApi("backend", "info", "service", "Received create short URL request for: "+req.URL)

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		logging.LogApi("backend", "error", "handler", "Invalid URL: "+req.URL)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	validity := req.Validity
	if validity <= 0 {
		validity = 30
	}

	shortcode := req.Shortcode
	if shortcode == "" {
		shortcode = generateShortcode()
	}

	urlStoreMutex.Lock()
	defer urlStoreMutex.Unlock()
	if _, exists := urlStore[shortcode]; exists {
		logging.LogApi("backend", "error", "handler", "Shortcode already exists: "+shortcode)
		http.Error(w, "Shortcode already exists", http.StatusConflict)
		return
	}

	expiry := time.Now().Add(time.Duration(validity) * time.Minute)
	shortURL := &ShortURL{
		OriginalURL: req.URL,
		Shortcode:   shortcode,
		CreatedAt:   time.Now(),
		Expiry:      expiry,
		Clicks:      0,
		ClickData:   []ClickDetail{},
	}
	urlStore[shortcode] = shortURL

	resp := CreateShortURLResponse{
		ShortLink: strings.TrimRight(r.Host, "/") + "/" + shortcode,
		Expiry:   expiry.UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	logging.LogApi("backend", "info", "service", "Short URL created: "+resp.ShortLink)
}

func getShortURLStatsHandler(w http.ResponseWriter, r *http.Request) {
	shortcode := chi.URLParam(r, "shortcode")
	urlStoreMutex.RLock()
	shortURL, exists := urlStore[shortcode]
	urlStoreMutex.RUnlock()
	if !exists {
		http.Error(w, "Shortcode not found", http.StatusNotFound)
		return
	}
	resp := ShortURLStatsResponse{
		Shortcode:   shortURL.Shortcode,
		OriginalURL: shortURL.OriginalURL,
		CreatedAt:   shortURL.CreatedAt.UTC().Format(time.RFC3339),
		Expiry:      shortURL.Expiry.UTC().Format(time.RFC3339),
		Clicks:      shortURL.Clicks,
		ClickData:   shortURL.ClickData,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortcode := chi.URLParam(r, "shortcode")
	urlStoreMutex.Lock()
	shortURL, exists := urlStore[shortcode]
	if !exists {
		urlStoreMutex.Unlock()
		http.Error(w, "Shortcode not found", http.StatusNotFound)
		return
	}
	if time.Now().After(shortURL.Expiry) {
		urlStoreMutex.Unlock()
		http.Error(w, "Short URL expired", http.StatusGone)
		return
	}
	// Record click
	shortURL.Clicks++
	click := ClickDetail{
		Timestamp: time.Now(),
		Referrer:  r.Referer(),
	}
	shortURL.ClickData = append(shortURL.ClickData, click)
	urlStoreMutex.Unlock()
	// Redirect
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusFound)
} 