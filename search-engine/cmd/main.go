package main

import (
	"search-engine/handlers" // Import the handler package

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Define the route to handle finding similar news
	r.GET("/find-similar-news", handlers.FindSimilarNews)

	r.Run(":8002")
}
