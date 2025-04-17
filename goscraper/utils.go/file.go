package utils

import (
	"encoding/csv"
	"os"
)

// SaveToCSV writes scraped data to a CSV file
func SaveToCSV(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(data)
}
