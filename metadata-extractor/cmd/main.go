package main

import (
	"log"
	"metadata-extractor/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/extract-metadata", handlers.GetMetadata)

	if err := r.Run(":8001"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
