package main

import (
	"content-extractor/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/get-content", handlers.GetContent)

	r.Run("localhost:8000")
}
