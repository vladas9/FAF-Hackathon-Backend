package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jpoz/groq"
)

type Message struct {
	Url  string `json:"url"`
	OUrl string `json:"original_url"`
}

func HandleGetSimilarity(c *gin.Context) {
	msg := &Message{}
	if err := c.ShouldBindJSON(msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content := calculateSimilarity(msg.OUrl, msg.Url)

	c.JSON(http.StatusOK, content)
}

func truncateText(text string, maxLength int) string {
	if len(text) > maxLength {
		return text[:maxLength] + "..." // Truncate and add ellipsis
	}
	return text
}

func calculateSimilarity(url1 string, url2 string) gin.H {
	godotenv.Load()

	text1, err := getContent(url1)
	text2, err := getContent(url2)

	if err != nil {
		fmt.Println(err)
	}

	maxLength := 1000

	truncatedText1 := truncateText(text1.Paragraphs, maxLength)
	truncatedText2 := truncateText(text2.Paragraphs, maxLength)

	client := groq.NewClient()
	var groqResponse *groq.ChatCompletion

	maxRetries := 3
	for attempts := 1; attempts <= maxRetries; attempts++ {
		groqResponse, err = client.CreateChatCompletion(groq.CompletionCreateParams{
			Model:          "gemma2-9b-it",
			ResponseFormat: groq.ResponseFormat{Type: "json_object"},
			Messages: []groq.Message{
				{
					Role:    "system",
					Content: "Always respond with valid JSON only. Do not include any additional text or explanation.",
				},
				{
					Role:    "user",
					Content: "Your task is to compare two news articles and provide a similarity score for each meaningful part based on the following criteria: Identify relevant sentences or segments that contain 10 or more words, Focus on comparing the factual content and information from both articles, Look for contradictions (e.g., conflicting facts), redundancies, and note how similar the ideas are in terms of meaning and context, Assign a similarity score between 0 and 100 for each part, categorizing them as 'Exact Match', 'Partial Match', or 'Contradiction'. If the articles are entirely unrelated, mark them as 'Not Related' with a similarity score of 0%,If the articles discuss different topics without significant overlap, indicate 'Not Related'. Please also identify similar ideas or phrasing, even if the wording differs.",
				},
				{
					Role: "user",
					Content: "Your answer should be in the following format:\n" +
						"{\n" +
						"    overallScore: ...," +
						"    similarParts: [" +
						"        { part1: 'Text from Article 1', part2: 'Text from Article 2', score: 'similarity score', typeOfSimilarity: 'Exact Match' | 'Partial Match' | 'Contradiction' }," +
						"        ..." +
						"    ]," +
						"    conflictingParts: [" +
						"        { part1: 'Text from Article 1', part2: 'Text from Article 2', score: 'similarity score', typeOfSimilarity: 'Exact Match' | 'Partial Match' | 'Contradiction' },\n" +
						"        ..." +
						"    ]," +
						"    notRelated: true // Indicate if the articles are unrelated, with a similarity score of 0%." +
						"}",
				},
				{
					Role:    "user",
					Content: fmt.Sprintf("Text1: %s, Text2: %s", truncatedText1, truncatedText2),
				},
			},
		})

		if err != nil {
			log.Printf("Error on attempt %d: %v", attempts, err)
		} else if groqResponse == nil || len(groqResponse.Choices) == 0 || len(groqResponse.Choices[0].Message.Content) == 0 {
			log.Printf("Empty response on attempt %d", attempts)
		} else {
			// Parse the JSON response
			var parsedResponse map[string]interface{}
			err = json.Unmarshal([]byte(groqResponse.Choices[0].Message.Content), &parsedResponse)
			if err != nil {
				log.Printf("Failed to parse JSON: %v", err)
				return gin.H{"error": "Failed to parse response JSON"}
			}
			return parsedResponse
		}

		if attempts < maxRetries {
			log.Printf("Retrying... attempt %d", attempts+1)
			time.Sleep(2 * time.Second)
		}
	}

	return gin.H{"error": "Failed to get response after multiple retries"}
}

// Define a struct for the parsed JSON
type SimilarityResponse struct {
	OverallScore int `json:"overallScore"`
	SimilarParts []struct {
		Part1            string `json:"part1"`
		Part2            string `json:"part2"`
		Score            int    `json:"score"`
		TypeOfSimilarity string `json:"typeOfSimilarity"`
	} `json:"similarParts"`
	ConflictingParts []struct {
		Part1            string `json:"part1"`
		Part2            string `json:"part2"`
		Score            int    `json:"score"`
		TypeOfSimilarity string `json:"typeOfSimilarity"`
	} `json:"conflictingParts"`
	NotRelated bool `json:"notRelated"`
}

func ParseInvalidJSON(input string) (*SimilarityResponse, error) {
	// Replace invalid placeholders with valid JSON nulls or empty strings
	cleanedInput := strings.ReplaceAll(input, `"%!s(*handlers.PageContent=<nil>)"`, `""`)

	// Unmarshal the cleaned JSON string into a struct
	var result SimilarityResponse
	err := json.Unmarshal([]byte(cleanedInput), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &result, nil
}

type PageContent struct {
	Paragraphs string   `json:"paragraphs"`
	Headings   []string `json:"headings"`
	Images     []string `json:"images"`
	Lists      []string `json:"lists"`
	Spans      []string `json:"spans"`
}

func getContent(url string) (*PageContent, error) {
	apiURL := "http://localhost:8000/get-content"
	jsonData := fmt.Sprintf("{\"url\": \"%s\"}", url)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var apiResponse *PageContent
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return apiResponse, nil
}

func getUrls(url string) {}
