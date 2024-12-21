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
	Paragraphs string   `json:"paragraphs"`
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

	paragraphLocator := page.Locator("p")
	paragraphCount, err := paragraphLocator.Count()
	if err != nil {
		log.Fatalf("could not get paragraph count: %v", err)
		return nil, err
	}

	var allParagraphs string
	for i := 0; i < paragraphCount; i++ {
		text, err := paragraphLocator.Nth(i).TextContent()
		if err != nil {
			log.Println("Error extracting paragraph text:", err)
			continue
		}
		allParagraphs += strings.TrimSpace(text) + " "
	}

	content.Paragraphs = strings.ReplaceAll(allParagraphs, `"`, `'`)

	for i := 1; i <= 6; i++ {
		headingTags := fmt.Sprintf("h%d", i)
		headingLocator := page.Locator(headingTags)
		headingCount, err := headingLocator.Count()
		if err != nil {
			log.Fatalf("could not get <%s> tag count: %v", headingTags, err)
			return nil, err
		}

		for j := 0; j < headingCount; j++ {
			text, err := headingLocator.Nth(j).TextContent()
			if err != nil {
				log.Println("Error extracting heading text:", err)
				continue
			}
			content.Headings = append(content.Headings, strings.TrimSpace(text))
		}
	}

	listLocator := page.Locator("ul li, ol li")
	listCount, err := listLocator.Count()
	if err != nil {
		log.Fatalf("could not get list item count: %v", err)
		return nil, err
	}

	for i := 0; i < listCount; i++ {
		text, err := listLocator.Nth(i).TextContent()
		if err != nil {
			log.Println("Error extracting list item text:", err)
			continue
		}
		content.Lists = append(content.Lists, strings.TrimSpace(text))
	}

	spanLocator := page.Locator("span")
	spanCount, err := spanLocator.Count()
	if err != nil {
		log.Fatalf("could not get span tag count: %v", err)
		return nil, err
	}

	for i := 0; i < spanCount; i++ {
		text, err := spanLocator.Nth(i).TextContent()
		if err != nil {
			log.Println("Error extracting span text:", err)
			continue
		}
		content.Spans = append(content.Spans, strings.TrimSpace(text))
	}

	imgLocator := page.Locator("img")
	imgCount, err := imgLocator.Count()
	if err != nil {
		log.Fatalf("could not get image count: %v", err)
		return nil, err
	}

	for i := 0; i < imgCount; i++ {
		src, err := imgLocator.Nth(i).GetAttribute("src")
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
