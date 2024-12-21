package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

type Message struct {
	Url string `json:"url"`
}

type Metadata struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	Title  string `json:"title"`
}

func GetMetadata(c *gin.Context) {
	msg := &Message{}
	if err := c.ShouldBindJSON(msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("URL: ", msg.Url)
	metadata, err := ExtractMetadataFromPage(msg.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

func ExtractMetadataFromPage(url string) (*Metadata, error) {
	metadata := &Metadata{}

	// Create a new collector
	c := colly.NewCollector()

	// Extract title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		metadata.Title = strings.TrimSpace(e.Text)
	})

	// Extract author (meta tag with name="author")
	c.OnHTML("meta[name='author']", func(e *colly.HTMLElement) {
		metadata.Author = strings.TrimSpace(e.Attr("content"))
	})

	// Extract date (meta tag with name="date")
	c.OnHTML("meta[name='date']", func(e *colly.HTMLElement) {
		metadata.Date = strings.TrimSpace(e.Attr("content"))
	})

	// Fallback: If no date in meta, try extracting from <time> tag
	c.OnHTML("time", func(e *colly.HTMLElement) {
		if metadata.Date == "" {
			metadata.Date = strings.TrimSpace(e.Text)
		}
	})

	// Start scraping the page
	err := c.Visit(url)
	if err != nil {
		log.Printf("Error visiting URL: %v", err)
		return nil, err
	}

	return metadata, nil
}
