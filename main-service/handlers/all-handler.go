package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleAll(c *gin.Context) {
	url := "https://point.md/ru/novosti/ekonomika/cheban-my-ozhidaem-uvedomleniia-ot-gazproma-v-blizhaishie-dni/"
	urls := getUrls(url)
	similar, _ := getSimilarity(urls, url)
	c.JSON(http.StatusOK, similar)
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

type Similar struct {
	OverallScore     int     `json:"overallScore"`
	SimilarParts     []Parts `json:"similarParts"`
	ConflictingParts []Parts `json:"conflictingParts"`
}

type Parts struct {
	Part1 string `json:"part1"`
	Part2 string `json:"part2"`
}

func getSimilarity(urls SimilarUrls, originalUrl string) (Similar, error) {

	apiURL := "http://localhost:8003/get-similarity"
	jsonData := fmt.Sprintf("{\"url\": \"%s\", \"original_url\": \"%s\"}", urls.Urls[2], originalUrl)

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

	var apiResp Similar
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		// return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return apiResp, nil

}

type Metadata struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	Title  string `json:"title"`
	Icon   string `json:"icon"`
}
