package scraper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// ParseTitle extracts the <title> from an HTML document
func ParseTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("title").Text()
	if title == "" {
		return "", fmt.Errorf("no title found")
	}
	return title, nil
}
