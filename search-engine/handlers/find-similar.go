package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type NewsMetadata struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

type DuckDuckGoResponse struct {
	Results []struct {
		URL string `json:"url"`
	} `json:"results"`
}

func parseDuckDuckGoURL(redirectURL string) (string, error) {
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	uddgParam := parsedURL.Query().Get("uddg")
	if uddgParam == "" {
		return "", fmt.Errorf("no 'uddg' parameter found in the URL")
	}

	decodedURL, err := url.QueryUnescape(uddgParam)
	if err != nil {
		return "", fmt.Errorf("failed to decode 'uddg' parameter: %v", err)
	}

	return decodedURL, nil
}

func FindSimilarNews(c *gin.Context) {
	url := c.DefaultQuery("url", "")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	metadataServiceURL := "http://localhost:8001/extract-metadata"

	payload := map[string]string{"url": url}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing request"})
		return
	}

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

	var metadata NewsMetadata
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		log.Printf("Error decoding metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding metadata"})
		return
	}

	log.Printf("Received metadata: %+v", metadata)

	searchQuery := metadata.Title
	duckDuckGoURL := fmt.Sprintf("https://duckduckgo.com/html/?q=%s", strings.ReplaceAll(searchQuery, " ", "+"))

	log.Printf("Constructed DuckDuckGo URL: %s", duckDuckGoURL)

	respDDG, err := http.Get(duckDuckGoURL)
	if err != nil {
		log.Printf("Error calling DuckDuckGo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calling DuckDuckGo"})
		return
	}
	defer respDDG.Body.Close()

	doc, err := goquery.NewDocumentFromReader(respDDG.Body)
	if err != nil {
		log.Printf("Error parsing DuckDuckGo response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing DuckDuckGo response"})
		return
	}

	var urlsToReturn []string
	doc.Find(".result__a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			decodedURL, err := parseDuckDuckGoURL(url)
			if err == nil {
				urlsToReturn = append(urlsToReturn, decodedURL)
				if len(urlsToReturn) >= 3 {
					return
				}
			}
		}
	})

	if len(urlsToReturn) == 0 {
		log.Println("No similar news found.")
		c.JSON(http.StatusOK, gin.H{
			"message": "No similar news found",
			"urls":    []string{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched similar news",
		"urls":    urlsToReturn,
	})
}
