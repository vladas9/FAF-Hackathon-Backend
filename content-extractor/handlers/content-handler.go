package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Message struct {
	Url string `json:"url"`
}

type Content struct {
	H   []string `json:"h"`
	P   []string `json:"p"`
	Div []string `json:"div"`
}

func GetContent(c *gin.Context) {
	msg := &Message{}
	if err := c.ShouldBindJSON(msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("URL: ", msg.Url)
	content := ExtractContentFromPage(msg.Url)

	c.JSON(http.StatusOK, content)
}

func ExtractContentFromPage(url string) []string {
	c := ""

	collector := colly.NewCollector()

	collector.OnHTML("body", func(e *colly.HTMLElement) {
		c = e.Text
	})

	collector.Visit(url)

	re := regexp.MustCompile(`"([^"]*)"`)
	matches := re.FindAllStringSubmatch(c, -1)

	var result []string

	for _, match := range matches {
		words := strings.Fields(match[1])
		if len(words) > 10 {
			result = append(result, match[1])
		}
	}

	return result
}
