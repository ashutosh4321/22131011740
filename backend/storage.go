package main

import (
	"sync"
	"time"
)

type ClickDetail struct {
	Timestamp time.Time
	Referrer  string
}

type ShortURL struct {
	OriginalURL string
	Shortcode   string
	CreatedAt   time.Time
	Expiry      time.Time
	Clicks      int
	ClickData   []ClickDetail
}

var urlStore = make(map[string]*ShortURL) // key: shortcode
var urlStoreMutex sync.RWMutex 