package scraper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// FetchTitle fetches and extracts the page title
func FetchTitle(url string) (string, error) {
	resp, err := FetchHTML(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Load HTML into goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract Title
	title, err := ParseTitle(doc)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Title of %s: %s", url, title), nil
}
