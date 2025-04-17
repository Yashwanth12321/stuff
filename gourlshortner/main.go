package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

const (
	DB_CONN  = "postgres://clickerdata:mysecretpassword@localhost:5432/urlshortener?sslmode=disable"
	BASE_URL = "http://localhost:8080/"
)

type URL struct {
	ID       int    `json:"id"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
	Clicks   int    `json:"clicks"`
}

var (
	db    *sql.DB
	rdb   *redis.Client
	ctx   = context.Background()
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func initDB() {
	var err error
	db, err = sql.Open("postgres", DB_CONN)
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_url VARCHAR(10) UNIQUE,
		long_url TEXT NOT NULL,
		clicks INT DEFAULT 0
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}
	fmt.Println("Database initialized successfully.")
}

func init() {
	initDB()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func generateShortURL() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func shortenURL(c *gin.Context) {
	var req struct {
		LongURL string `json:"long_url"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	shortURL := generateShortURL()
	_, err := db.Exec("INSERT INTO urls (short_url, long_url, clicks) VALUES ($1, $2, $3)", shortURL, req.LongURL, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rdb.Set(ctx, shortURL, req.LongURL, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"short_url": BASE_URL + shortURL})
}

func redirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	longURL, err := rdb.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		err = db.QueryRow("SELECT long_url FROM urls WHERE short_url = $1", shortURL).Scan(&longURL)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		rdb.Set(ctx, shortURL, longURL, 10*time.Minute)
	}

	_, _ = db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short_url = $1", shortURL)
	c.Redirect(http.StatusFound, longURL)
}

func getAnalytics(c *gin.Context) {
	rows, err := db.Query("SELECT short_url, long_url, clicks FROM urls ORDER BY clicks DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var results []URL
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.ShortURL, &url.LongURL, &url.Clicks); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data parsing error"})
			return
		}
		results = append(results, url)
	}

	c.JSON(http.StatusOK, results)
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/shorten", shortenURL)
	r.GET("/:shortURL", redirectURL)
	r.GET("/analytics", getAnalytics)

	r.Run(":8080")
}
