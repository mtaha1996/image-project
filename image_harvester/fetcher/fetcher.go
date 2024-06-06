package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mtaha1996/image-project/image_harvester/config"
)

type SearchResult struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

func FetchImageURLs(cfg config.Config, query string, maxImages int) ([]string, error) {
	encodedQuery := url.QueryEscape(query)

	searchURL := fmt.Sprintf(
		"%s/customsearch/v1?q=%s&key=%s&cx=%s&searchType=image&num=%d",
		cfg.BaseURL, encodedQuery, cfg.GoogleAPIKey, cfg.SearchEngineID, maxImages)

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
