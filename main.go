package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	// ping
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// static files
	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")

	router.Run(":8080")
}
