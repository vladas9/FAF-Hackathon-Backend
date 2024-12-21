package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/playwright-community/playwright-go"
)

type Message struct {
	Url string `json:"url"`
}

type Content struct {
	Paragraphs []string `json:"paragraphs"`
	Headings   []string `json:"headings"`
	Images     []string `json:"images"`
	Lists      []string `json:"lists"`
	Spans      []string `json:"spans"`
}

func GetContent(c *gin.Context) {
	msg := &Message{}
	if err := c.ShouldBindJSON(msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("URL: ", msg.Url)
	content, err := ExtractContentFromPage(msg.Url)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, content)
}

func ExtractContentFromPage(url string) (*Content, error) {
	content := &Content{}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
		return nil, err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
		return nil, err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
		return nil, err
	}

	_, err = page.Goto(url)
	if err != nil {
		log.Fatalf("could not visit page: %v", err)
		return nil, err
	}

	_, err = page.WaitForSelector("body")
	if err != nil {
		log.Fatalf("could not wait for selector: %v", err)
		return nil, err
	}

	paragraphs, err := page.QuerySelectorAll("p")
	if err != nil {
		log.Fatalf("could not get <p> tags: %v", err)
		return nil, err
	}

	for _, p := range paragraphs {
		text, err := p.TextContent()
		if err != nil {
			log.Println("Error extracting paragraph text:", err)
			continue
		}
		content.Paragraphs = append(content.Paragraphs, strings.TrimSpace(text))
	}

	for i := 1; i <= 6; i++ {
		headingTags := fmt.Sprintf("h%d", i)
		headings, err := page.QuerySelectorAll(headingTags)
		if err != nil {
			log.Fatalf("could not get <%s> tags: %v", headingTags, err)
			return nil, err
		}

		for _, h := range headings {
			text, err := h.TextContent()
			if err != nil {
				log.Println("Error extracting heading text:", err)
				continue
			}
			content.Headings = append(content.Headings, strings.TrimSpace(text))
		}
	}

	listItems, err := page.QuerySelectorAll("ul li, ol li")
	if err != nil {
		log.Fatalf("could not get list items: %v", err)
		return nil, err
	}

	for _, li := range listItems {
		text, err := li.TextContent()
		if err != nil {
			log.Println("Error extracting list item text:", err)
			continue
		}
		content.Lists = append(content.Lists, strings.TrimSpace(text))
	}

	spans, err := page.QuerySelectorAll("span")
	if err != nil {
		log.Fatalf("could not get <span> tags: %v", err)
		return nil, err
	}

	for _, span := range spans {
		text, err := span.TextContent()
		if err != nil {
			log.Println("Error extracting span text:", err)
			continue
		}
		content.Spans = append(content.Spans, strings.TrimSpace(text))
	}

	images, err := page.QuerySelectorAll("img")
	if err != nil {
		log.Fatalf("could not get <img> tags: %v", err)
		return nil, err
	}

	for _, img := range images {
		src, err := img.GetAttribute("src")
		if err != nil {
			log.Println("Error extracting image src:", err)
			continue
		}
		if src != "" {
			content.Images = append(content.Images, src)
		}
	}

	return content, nil
}
