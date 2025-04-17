package config

import "time"

var Config = struct {
	Timeout   time.Duration
	UserAgent string
}{
	Timeout:   10 * time.Second,
	UserAgent: "Mozilla/5.0 (compatible; GoScraper/1.0)",
}
