package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewsMetadata represents the metadata structure expected from the metadata-extractor microservice
type NewsMetadata struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

// NewsAPIResponse represents the structure of the response from NewsAPI
type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		URL string `json:"url"`
	} `json:"articles"`
}

// FindSimilarNews handles the incoming request, fetches metadata from the metadata-extractor service,
// queries NewsAPI for similar news, and returns URLs of matching articles
func FindSimilarNews(c *gin.Context) {
	// Get the URL from the query parameter
	url := c.DefaultQuery("url", "")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// URL of the metadata-extractor microservice
	metadataServiceURL := "http://localhost:8001/extract-metadata"

	// Create the request payload
	payload := map[string]string{"url": url}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing request"})
		return
	}

	// Send POST request to metadata-extractor
	resp, err := http.Post(metadataServiceURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error fetching metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching metadata"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code from metadata service: %d", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected response from metadata service"})
		return
	}

	// Parse the metadata response
	var metadata NewsMetadata
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		log.Printf("Error decoding metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding metadata"})
		return
	}

	// Log the metadata for debugging
	log.Printf("Received metadata: %+v", metadata)

	// Construct the NewsAPI query
	searchQuery := metadata.Title
	newsAPIKey := "your_api_key_here" // Replace with your actual NewsAPI key
	newsAPIURL := fmt.Sprintf(
		"https://newsapi.org/v2/everything?q=%s&from=%s&apiKey=%s",
		strings.ReplaceAll(searchQuery, " ", "+"),
		metadata.Date,
		newsAPIKey,
	)

	// Log the constructed URL for debugging
	log.Printf("Constructed NewsAPI URL: %s", newsAPIURL)

	// Call NewsAPI
	respAPI, err := http.Get(newsAPIURL)
	if err != nil {
		log.Printf("Error calling NewsAPI: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calling NewsAPI"})
		return
	}
	defer respAPI.Body.Close()

	// Parse the NewsAPI response
	var newsResponse NewsAPIResponse
	err = json.NewDecoder(respAPI.Body).Decode(&newsResponse)
	if err != nil {
		log.Printf("Error decoding NewsAPI response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding NewsAPI response"})
		return
	}

	// Log the NewsAPI response for debugging
	log.Printf("NewsAPI Response: %+v", newsResponse)

	if newsResponse.TotalResults == 0 {
		log.Println("No similar news found.")
		c.JSON(http.StatusOK, gin.H{
			"message": "No similar news found",
			"urls":    []string{},
		})
		return
	}

	// Extract URLs from the response
	var urls []string
	for _, article := range newsResponse.Articles {
		urls = append(urls, article.URL)
	}

	// Return the results
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched similar news",
		"urls":    urls,
	})
}
