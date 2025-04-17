package utils

import (
	"log"
	"os"
)

// Initialize logger
var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
