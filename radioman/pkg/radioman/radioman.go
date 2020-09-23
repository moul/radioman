package radioman

import (
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func main2() {
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
	public := router.Group("/")
	admin := router.Group("/")
	liquidsoap := router.Group("/")
	// Admin auth
	// FIXME: make accounts dynamic
	accounts := gin.Accounts{"admin": "admin"}
	admin.Use(gin.BasicAuth(accounts))
	// FIXME: add authentication on liquidsoap next handler

	public.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	staticPrefix := "./radioman/web"
	if os.Getenv("WEBDIR") != "" {
		staticPrefix = os.Getenv("WEBDIR")
	}

	public.StaticFile("/", path.Join(staticPrefix, "static/index.html"))
	public.Static("/static", path.Join(staticPrefix, "static"))
	public.Static("/bower_components", path.Join(staticPrefix, "bower_components"))

	admin.StaticFile("/admin/", path.Join(staticPrefix, "static/admin/index.html"))

	admin.GET("/api/playlists", playlistsEndpoint)
	admin.GET("/api/playlists/:name", playlistDetailEndpoint)
	admin.PATCH("/api/playlists/:name", playlistUpdateEndpoint)
	admin.GET("/api/playlists/:name/tracks", playlistTracksEndpoint)

	admin.GET("/api/radios/default", defaultRadioEndpoint)
	public.GET("/api/radios/default/endpoints", radioEndpointsEndpoint)
	admin.POST("/api/radios/default/skip-song", radioSkipSongEndpoint)
	admin.POST("/api/radios/default/play-track", radioPlayTrackEndpoint)
	admin.POST("/api/radios/default/set-next-track", radioSetNextTrackEndpoint)

	admin.GET("/api/tracks/:hash", trackDetailEndpoint)

	liquidsoap.GET("/api/liquidsoap/getNextSong", getNextSongEndpoint)

	public.GET("/playlist.m3u", m3uPlaylistEndpoint)

	// Launch routines
	go Radio.UpdatePlaylistsRoutine()

	// Start web server mainloop
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(fmt.Sprintf(":%s", port))
}

type Opts struct {
	Verbose bool
}
