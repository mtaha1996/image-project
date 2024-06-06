package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/google/uuid"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nfnt/resize"
)

const (
	googleAPIKey            = "GOOGLE_API_KEY"
	searchEngineID          = "SEARCH_ENGINE_ID"
	baseUrl                 = "BASE_URL"
	IMAGE_POSTGRES_DB       = "IMAGE_POSTGRES_DB"
	IMAGE_POSTGRES_USER     = "IMAGE_POSTGRES_USER"
	IMAGE_POSTGRES_PASSWORD = "IMAGE_POSTGRES_PASSWORD"
	IMAGE_POSTGRES_HOST     = "IMAGE_POSTGRES_HOST"
	IMAGE_POSTGRES_PORT     = "IMAGE_POSTGRES_PORT"
)

type SearchResult struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

func fetchImageURLs(query string, maxImages int) ([]string, error) {
	apiKey := os.Getenv(googleAPIKey)
	cx := os.Getenv(searchEngineID)
	baseURL := os.Getenv(baseUrl)

	encodedQuery := url.QueryEscape(query)

	searchURL := fmt.Sprintf(
		"%s/customsearch/v1?q=%s&key=%s&cx=%s&searchType=image&num=%d",
		baseURL, encodedQuery, apiKey, cx, maxImages)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResult SearchResult
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, err
	}

	var imageUrls []string
	for _, item := range searchResult.Items {
		imageUrls = append(imageUrls, item.Link)
	}

	return imageUrls, nil
}

func downloadAndResizeImage(url string, size uint) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	m := resize.Resize(size, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, m, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func saveImageToDB(db *sql.DB, filename, query string, imageData []byte) error {
	_, err := db.Exec("INSERT INTO images(file_name, query, image_data) VALUES($1, $2, $3)", filename, query, imageData)
	return err
}

func saveImageToFile(filename, directory string, imageData []byte) error {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Errorf("failed to create directory: %v", err)
	}

	// Construct the full path for the file
	filePath := filepath.Join(directory, filename)

	// Write the image data to the file
	err = os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		fmt.Errorf("failed to save image to file: %v", err)
	}

	return nil
}

func checkAndCreateTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS images (
		id SERIAL PRIMARY KEY,
		query TEXT,
		file_name TEXT,
		image_data BYTEA
	)`
	_, err := db.Exec(createTableQuery)
	return err
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <query> <max_images> <image_width>\n", os.Args[0])
	}

	query := os.Args[1]

	maxImages, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid max_images value: %s", os.Args[2])
	}
	imageWidth, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Invalid imageWidth value: %s", os.Args[2])
	}

	dbHost := os.Getenv(IMAGE_POSTGRES_HOST)
	dbPort, _ := strconv.Atoi(os.Getenv(IMAGE_POSTGRES_PORT))
	dbUser := os.Getenv(IMAGE_POSTGRES_USER)
	dbPassword := os.Getenv(IMAGE_POSTGRES_PASSWORD)
	dbName := os.Getenv(IMAGE_POSTGRES_DB)

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		" password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	err = checkAndCreateTable(db)
	if err != nil {
		log.Fatalf("Failed to check or create table: %v", err)
	}

	imageURLs, err := fetchImageURLs(query, maxImages)
	if err != nil {
		log.Fatalf("Failed to fetch image URLs: %v", err)
	}

	var wg sync.WaitGroup
	for _, url := range imageURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			imageData, err := downloadAndResizeImage(url, uint(imageWidth))
			if err != nil {
				log.Printf("Failed to download and resize image: %v", err)
				return
			}

			// Generate a random file name
			filename := uuid.New().String() + ".jpg"

			err = saveImageToDB(db, filename, query, imageData)
			if err != nil {
				log.Printf("Failed to save image to database: %v", err)
			}

			// Save image to directory with a random file name
			directory := "./downloaded_images"
			err = saveImageToFile(filename, directory, imageData)
			if err != nil {
				log.Printf("Failed to save image to directory: %v", err)
			} else {
				log.Printf("Image successfully saved to directory with filename: %s", filename)
			}

		}(url)
	}
	wg.Wait()

	log.Println("Done!")
}
