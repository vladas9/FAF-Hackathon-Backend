package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url" // Import the net/url package
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// NewsMetadata represents the metadata structure expected from the metadata-extractor microservice
type NewsMetadata struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

// DuckDuckGoResponse represents the structure of the response from DuckDuckGo search
type DuckDuckGoResponse struct {
	Results []struct {
		URL string `json:"url"`
	} `json:"results"`
}

// FindSimilarNews handles the incoming request, fetches metadata from the metadata-extractor service,
// queries DuckDuckGo for similar news, and returns a single URL of a matching article
func FindSimilarNews(c *gin.Context) {
	// Get the URL from the query parameter
	originalURL := c.DefaultQuery("url", "") // Renamed the variable to avoid conflict
	if originalURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// URL of the metadata-extractor microservice
	metadataServiceURL := "http://localhost:8001/extract-metadata"

	// Create the request payload
	payload := map[string]string{"url": originalURL}
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

	// Construct the DuckDuckGo search query
	searchQuery := metadata.Title
	duckDuckGoURL := fmt.Sprintf("https://duckduckgo.com/html/?q=%s", strings.ReplaceAll(searchQuery, " ", "+"))

	// Log the constructed URL for debugging
	log.Printf("Constructed DuckDuckGo URL: %s", duckDuckGoURL)

	// Call DuckDuckGo search
	respDDG, err := http.Get(duckDuckGoURL)
	if err != nil {
		log.Printf("Error calling DuckDuckGo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calling DuckDuckGo"})
		return
	}
	defer respDDG.Body.Close()

	// Parse the DuckDuckGo response using goquery
	doc, err := goquery.NewDocumentFromReader(respDDG.Body)
	if err != nil {
		log.Printf("Error parsing DuckDuckGo response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing DuckDuckGo response"})
		return
	}

	// Extract the first valid URL from the DuckDuckGo search results
	var urlToReturn string
	doc.Find(".result__a").Each(func(i int, s *goquery.Selection) {
		urlStr, exists := s.Attr("href")
		if exists {
			// Check if the URL is relative, and prepend the base URL if necessary
			if strings.HasPrefix(urlStr, "/url?q=") {
				// Extract the actual URL from the query parameter
				urlStr = strings.TrimPrefix(urlStr, "/url?q=")
				// Decode the URL if necessary
				decodedURL, err := url.QueryUnescape(urlStr) // Correctly using url.QueryUnescape
				if err != nil {
					log.Printf("Error decoding URL: %v", err)
					return
				}
				// Use the decoded URL
				urlStr = decodedURL
			}
			// Set the first valid URL and stop the loop
			urlToReturn = urlStr
			return
		}
	})

	// Check if a valid URL was found
	if urlToReturn == "" {
		log.Println("No similar news found.")
		c.JSON(http.StatusOK, gin.H{
			"message": "No similar news found",
			"url":     "",
		})
		return
	}

	// Return the first valid URL
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched similar news",
		"url":     urlToReturn,
	})
}
