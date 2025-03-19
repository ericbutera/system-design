package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// runtime.GOMAXPROCS(1)

	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "the time is " + time.Now().String()})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
