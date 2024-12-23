package main

import (
	"main-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/get-info", handlers.HandleAll)

	r.Run(":6969")
}
