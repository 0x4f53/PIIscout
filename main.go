package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	apiKey         = "YOUR_GOOGLE_API_KEY"
	searchEngineID = "YOUR_CUSTOM_SEARCH_ENGINE_ID"
)

type SearchResult struct {
	Items []struct {
		Title string `json:"title"`
		Link  string `json:"link"`
		Image struct {
			ContextLink   string `json:"contextLink"`
			ThumbnailLink string `json:"thumbnailLink"`
		} `json:"image"`
	} `json:"items"`
}

func searchImages(query string) (*SearchResult, error) {
	baseURL := "https://www.googleapis.com/customsearch/v1"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("cx", searchEngineID)
	params.Add("q", query)
	params.Add("searchType", "image")

	searchURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search: %v", resp.Status)
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a search query")
	}

	query := os.Args[1]
	result, err := searchImages(query)
	if err != nil {
		log.Fatalf("Failed to search images: %v", err)
	}

	for _, item := range result.Items {
		fmt.Printf("Title: %s\nLink: %s\nThumbnail: %s\n\n", item.Title, item.Link, item.Image.ThumbnailLink)
	}
}
