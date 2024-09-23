package main

import (
	"github.com/gin-gonic/gin"
	"github.com/neverloved-dev/goJWT/db"
)

func handleReturnPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	r := gin.Default()
	db.Connect()
	r.GET("/ping", handleReturnPong)
	r.Run(":9000")
}
