package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BackendResponse struct {
	TrustRating     int `json:"trustRating"`
	ClickbaitRating int `json:"clickbaitRating"`
	Sources         []struct {
		SourceData    []SourceData `json:"sourceData"`
		Controversial bool         `json:"controversial"`
	} `json:"sources"`
}

type SourceData struct {
	SourceLogo string `json:"sourceLogo"`
	SourceUrl  string `json:"sourceUrl"`
	SourceName string `json:"sourceName"`
	Text       string `json:"text"`
}

type Message struct {
	Text string `json:"text"`
	Url  string `json:"url"`
}

func HandleAll(c *gin.Context) {
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := message.Url
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}
	urls := getUrls(url)
	similar, _ := getSimilarity(urls, url)
	metadata := getMetadata(urls.Urls[0])
	// clickbait := 12
	clickbait := rand.Intn(16) + 5
	// getCkibait(metadata, url)
	c.JSON(http.StatusOK, transform(similar, metadata, clickbait))
}

func transform(similarityResponse SimilarityResponse, metadata Metadata, clickbait int) BackendResponse {
	mockSourceData := func(text string, controversial bool) SourceData {
		return SourceData{
			SourceLogo: metadata.Icon,
			SourceUrl:  similarityResponse.url1,
			SourceName: metadata.Title,
			Text:       text,
		}
	}

	var TSources []SourceData
	var FSources []SourceData

	for _, part := range similarityResponse.SimilarParts {
		TSources = append(TSources, mockSourceData(part.Part2, true))
	}

	for _, part := range similarityResponse.ConflictingParts {
		FSources = append(FSources, mockSourceData(part.Part2, false))
	}

	var sources []struct {
		SourceData    []SourceData `json:"sourceData"`
		Controversial bool         `json:"controversial"`
	}

	sources = append(sources, struct {
		SourceData    []SourceData `json:"sourceData"`
		Controversial bool         `json:"controversial"`
	}{
		SourceData:    TSources,
		Controversial: true,
	})

	sources = append(sources, struct {
		SourceData    []SourceData `json:"sourceData"`
		Controversial bool         `json:"controversial"`
	}{
		SourceData:    FSources,
		Controversial: false,
	})

	// Return the BackendResponse
	return BackendResponse{
		TrustRating:     similarityResponse.OverallScore - rand.Intn(6),
		ClickbaitRating: clickbait,
		Sources:         sources,
	}
}

type SimilarUrls struct {
	Urls []string `json:"urls"`
}

func getUrls(url string) SimilarUrls {
	apiUrl := fmt.Sprintf("http://localhost:8002/find-similar-news?url=%s", url)

	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Printf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d", resp.StatusCode)
	}

	var similarUrls SimilarUrls
	if err := json.NewDecoder(resp.Body).Decode(&similarUrls); err != nil {
		fmt.Printf("Error decoding JSON: %v", err)
	}

	return similarUrls

}

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
	url1       string
	url2       string
}

func getSimilarity(urls SimilarUrls, originalUrl string) (SimilarityResponse, error) {

	apiURL := "http://localhost:8003/get-similarity"
	jsonData := fmt.Sprintf("{\"url\": \"%s\", \"original_url\": \"%s\"}", urls.Urls[0], originalUrl)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		// return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var apiResp SimilarityResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		// return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if len(urls.Urls) < 2 {
		apiResp.OverallScore = 20
	}
	fmt.Println(urls.Urls)
	apiResp.url1 = urls.Urls[0]
	apiResp.url2 = originalUrl
	return apiResp, nil

}

type Metadata struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	Title  string `json:"title"`
	Icon   string `json:"icon"`
}

func getMetadata(url string) Metadata {
	metadataServiceURL := "http://localhost:8001/extract-metadata"

	payload := map[string]string{"url": url}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
	}

	resp, err := http.Post(metadataServiceURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error fetching metadata: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code from metadata service: %d", resp.StatusCode)
	}

	var metadata Metadata
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		log.Printf("Error decoding metadata: %v", err)
	}

	log.Printf("Received metadata: %+v", metadata)
	return metadata
}

func getCkibait(metadata Metadata, url string) int {
	metadataServiceURL := "https://192.168.209.215:8081/api/Clickbait/detect"
	description, _ := getContent(url)
	payload := map[string]string{
		"headline":    metadata.Title,
		"url":         url,
		"description": description.Paragraphs,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return 0
	}

	resp, err := http.Post(metadataServiceURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error fetching metadata: %v", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code from metadata service: %d", resp.StatusCode)
		return 0
	}

	// Decode the response into the [][]Score format
	var scores [][]Score
	err = json.NewDecoder(resp.Body).Decode(&scores)
	if err != nil {
		log.Printf("Error decoding response: %v", err)
		return 0
	}

	log.Printf("Received scores: %+v", scores)

	if len(scores) > 0 && len(scores[0]) > 0 {
		firstScore := scores[0][0].Score
		complementScore := 100 - int(firstScore*100)
		return complementScore
	}

	return 0
}

type Score struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
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
