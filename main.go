package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
	"log"
)

const apiURL = "https://www.googleapis.com/customsearch/v1"

// SearchResult holds the structure of a single search result
type SearchResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// SearchResponse holds the structure of the API response
type SearchResponse struct {
	Items []SearchResult `json:"items"`
}

func FetchPIILinks(apiKey, cx, query, from, to string, includeLabel bool, outputFile, folderPath string) error {
	log.Println("Fetching links for query:", query) // Debug print
	dateRange, err := CalculateDateRange(from, to)
	if err != nil {
		return err
	}

	results, err := SearchGoogle(apiKey, cx, query, dateRange)
	if err != nil {
		return err
	}

	log.Println("Fetched results:", len(results.Items)) // Debug print to see if results are fetched

	output := FormatResults(results)

	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}

	// Write the results in JSON format (one result per line)
	if outputFile != "" {
		// Open the file in append mode, create it if it doesn't exist
		file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		if _, err := file.WriteString(output); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	} else {
		log.Println(output) // Print the results to console
	}
	return nil
}

func SearchGoogle(apiKey, cx, query, dateRange string) (*SearchResponse, error) {
	fmt.Printf("Searching Google with query: %s\n", query) // Debug print
	reqURL := fmt.Sprintf("%s?key=%s&cx=%s&q=%s&dateRestrict=%s", apiURL, apiKey, cx, url.QueryEscape(query), dateRange)

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	return &searchResponse, nil
}

// CalculateDateRange parses a date range in the format yyyy-mm-dd-hh
func CalculateDateRange(from, to string) (string, error) {
	const layout = "2006-01-02-15"
	var startDate, endDate time.Time
	var err error

	if from == "" && to == "" {
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -7) // Default to last week
	} else {
		startDate, err = time.Parse(layout, from)
		if err != nil {
			return "", fmt.Errorf("invalid 'from' date format, expected yyyy-mm-dd-hh")
		}
		endDate, err = time.Parse(layout, to)
		if err != nil {
			return "", fmt.Errorf("invalid 'to' date format, expected yyyy-mm-dd-hh")
		}
	}

	return fmt.Sprintf("%s,%s", startDate.Format("2006-01-02-15"), endDate.Format("2006-01-02-15")), nil
}

// FormatResults formats results in JSON format, one result per line
func FormatResults(results *SearchResponse) string {
	var output string
	for _, item := range results.Items {
		resultJSON, err := json.Marshal(item)
		if err != nil {
			fmt.Printf("Error marshalling result: %v\n", err)
			continue
		}
		output += string(resultJSON) + "\n" // Each result on a new line
	}
	return output
}

// DownloadFile downloads a file from the given URL
func DownloadFile(fileURL, folderPath string, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done() // Decrement the counter when the function completes

	resp, err := http.Get(fileURL)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch URL: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("failed to download file, status: %v", resp.Status)
		return
	}

	fileName := filepath.Base(fileURL)
	filePath := filepath.Join(folderPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		errChan <- fmt.Errorf("failed to create file: %v", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		errChan <- fmt.Errorf("failed to save file: %v", err)
		return
	}

	fmt.Printf("Downloaded: %s\n", fileName)
}

func LoadQueriesFromYAML(filePath string) ([]string, error) {
	var data struct {
		Dorks []string `yaml:"dorks"`
	}

	fileContent, err := os.ReadFile(filePath) // Replaced ioutil.ReadFile with os.ReadFile
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if err := yaml.Unmarshal(fileContent, &data); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	return data.Dorks, nil
}

func main() {
	// Command-line arguments
	var (
		apiKey   string
		cx       string
		country  string
		fromDate string
		toDate   string
	)

	flag.StringVar(&apiKey, "apiKey", "", "Google API key")
	flag.StringVar(&cx, "cx", "", "Custom search engine ID")
	flag.StringVar(&country, "country", "", "Country name (e.g., canada, india, pakistan)")
	flag.StringVar(&fromDate, "from", "", "Start date (yyyy-mm-dd-hh)")
	flag.StringVar(&toDate, "to", "", "End date (yyyy-mm-dd-hh)")
	flag.Parse()

	// Set includeLabel to true by default
	includeLabel := true

	if country == "" {
		log.Println("Error: --country is required")
		os.Exit(1)
	}

	// Generate output file name with country and timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	outputFile := fmt.Sprintf("%s-%s.json", country, timestamp)

	// Create the output folder with country name
	folderPath := filepath.Join("downloads", country, timestamp)

	// Verify and load country-specific queries
	dorkFile := filepath.Join("dorks", fmt.Sprintf("%s.yaml", country))
	if _, err := os.Stat(dorkFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Error: Country '%s' not found\n", country)
		os.Exit(1)
	}

	queries, err := LoadQueriesFromYAML(dorkFile)
	if err != nil {
		fmt.Printf("Error loading dork file: %v\n", err)
		os.Exit(1)
	}

	// Iterate through queries and execute
	for _, query := range queries {
		fmt.Printf("Executing query: %s\n", query)
		if err := FetchPIILinks(apiKey, cx, query, fromDate, toDate, includeLabel, outputFile, folderPath); err != nil {
			fmt.Printf("Error executing query: %v\n", err)
		}
	}

	// After generating the JSON output, download the files based on the links
	file, err := os.Open(outputFile)
	if err != nil {
		fmt.Printf("Error opening output file for downloading: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var result SearchResult
	decoder := json.NewDecoder(file)

	// Set up goroutines for concurrent downloads
	var wg sync.WaitGroup
	errChan := make(chan error, 10) // Buffered channel to collect errors

	// Loop through the results and download concurrently
	for decoder.More() {
		if err := decoder.Decode(&result); err != nil {
			fmt.Printf("Error decoding JSON result: %v\n", err)
			continue
		}

		wg.Add(1)
		go DownloadFile(result.Link, folderPath, &wg, errChan)
	}

	// Wait for all downloads to complete
	wg.Wait()

	// Check for any errors
	close(errChan)
	for err := range errChan {
		fmt.Printf("Download error: %v\n", err)
	}
}
