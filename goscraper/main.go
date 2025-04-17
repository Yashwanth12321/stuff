package main

import (
	"fmt"
	"sync"
	"webscraper/scraper"
)

func main() {
	urls := []string{
		"https://invalid.com",
		"https://example.com",
	}

	var wg sync.WaitGroup
	ch := make(chan string, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			title, err := scraper.FetchTitle(u)
			if err != nil {
				ch <- fmt.Sprintf("Error: %v", err)
			} else {
				ch <- title
			}
		}(url)
	}

	wg.Wait()
	close(ch)

	for msg := range ch {
		fmt.Println(msg)
	}
}
