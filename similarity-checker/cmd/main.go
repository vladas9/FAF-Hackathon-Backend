package main

import (
	"similarity-checker/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/get-similarity", handlers.HandleGetSimilarity)

	r.Run(":8003")
}
