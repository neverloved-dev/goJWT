package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func handleReturnPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	fmt.Println("hello world")
	r := gin.Default()
	r.GET("/ping", handleReturnPong)
	r.Run(":9000")
}
