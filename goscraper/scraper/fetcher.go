package scraper

import (
	"fmt"
	"net/http"

	"webscraper/config"
)

// FetchHTML retrieves the raw HTML from a URL
func FetchHTML(url string) (*http.Response, error) {
	client := &http.Client{Timeout: config.Config.Timeout}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", config.Config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, url)
	}
	return resp, nil
}
