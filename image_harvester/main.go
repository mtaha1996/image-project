package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"

	"github.com/image-downloader/image_harvester/config"
	"github.com/image-downloader/image_harvester/downloader"
	"github.com/image-downloader/image_harvester/fetcher"
	"github.com/image-downloader/image_harvester/storage"
)

func main() {
	cfg := config.LoadConfig()

	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <search_query> <max_images> <image_width>", os.Args[0])
	}

	query := os.Args[1]
	maxImages, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid maxImages value: %s", os.Args[2])
	}
	imageWidth, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Invalid imageWidth value: %s", os.Args[2])
	}

	db, err := storage.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = storage.CheckAndCreateTable(db)
	if err != nil {
		log.Fatalf("Failed to check or create table: %v", err)
	}

	imageURLs, err := fetcher.FetchImageURLs(cfg, query, maxImages)
	if err != nil {
		log.Fatalf("Failed to fetch image URLs: %v", err)
	}

	var wg sync.WaitGroup
	for _, url := range imageURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			imageData, err := downloader.DownloadAndResizeImage(url, uint(imageWidth))
			if err != nil {
				log.Printf("Failed to download and resize image: %v", err)
				return
			}

			filename := uuid.New().String() + ".jpg"

			err = storage.SaveImageToDB(db, filename, query, imageData)
			if err != nil {
				log.Printf("Failed to save image to database: %v", err)
			}

			err = downloader.SaveImageToFile(filename, cfg.ImageDirectoryName, imageData)
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
