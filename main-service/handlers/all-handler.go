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
	urls := getUrls("https://fortune.com/2024/12/21/government-shutdown-congress-spending-bill-donald-trump-debt-ceiling-elon-musk/")
	similar, _ := getSimilarity(urls)
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
	OverallScore     int   `json:"overallScore"`
	SimilarParts     Parts `json:"similarParts"`
	ConflictingParts Parts `json:"conflictingParts"`
}

type Parts struct {
	Part1 string `json:"part1"`
	Part2 string `json:"part2"`
}

func getSimilarity(urls SimilarUrls) ([]Similar, error) {

	var apiResponse []Similar

	fmt.Println("url: ", urls)
	for _, url := range urls.Urls {
		apiURL := "http://localhost:8003/get-similarity"
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

		var apiResp Similar
		err = json.Unmarshal(body, &apiResp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %v", err)
		}

		apiResponse = append(apiResponse, apiResp)
	}

	return apiResponse, nil

}

type Metadata struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	Title  string `json:"title"`
	Icon   string `json:"icon"`
}
