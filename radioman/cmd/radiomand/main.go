package main

import (
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/moul/radioman/radioman/pkg/radioman"
)

var Radio *radioman.Radio

func main() {
	// Setup the Radio instance
	Radio = radioman.NewRadio("RadioMan")

	if err := Radio.Init(); err != nil {
		logrus.Fatalf("Failed to initialize the radio: %v", err)
	}
	if err := Radio.StdPopulate(); err != nil {
		logrus.Fatalf("Failed to populate the radio: %v", err)
	}

	// Setup the web server
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	staticPrefix := "./radioman/web"
	if os.Getenv("WEBDIR") != "" {
		staticPrefix = os.Getenv("WEBDIR")
	}
	router.StaticFile("/", path.Join(staticPrefix, "static/index.html"))
	router.Static("/static", path.Join(staticPrefix, "static"))
	router.Static("/bower_components", path.Join(staticPrefix, "bower_components"))

	router.GET("/api/playlists", playlistsEndpoint)
	router.GET("/api/playlists/:name", playlistDetailEndpoint)
	router.PATCH("/api/playlists/:name", playlistUpdateEndpoint)
	router.GET("/api/playlists/:name/tracks", playlistTracksEndpoint)

	router.GET("/api/radios/default", defaultRadioEndpoint)
	router.POST("/api/radios/default/skip-song", radioSkipSongEndpoint)

	router.GET("/api/tracks/:hash", trackDetailEndpoint)

	router.GET("/api/liquidsoap/getNextSong", getNextSongEndpoint)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Launch routines
	go Radio.UpdatePlaylistsRoutine()

	// Start web server mainloop
	router.Run(fmt.Sprintf(":%s", port))
}
