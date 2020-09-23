package radioman

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	chilogger "github.com/treastech/logger"
)

func (r *Radio) server() *http.Server {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(cors.Handler)
	router.Use(chilogger.Logger(r.logger.Named("http")))
	router.Use(middleware.Timeout(time.Second * 5))
	router.Use(middleware.Recoverer)

	/*
		r.Route("/api", func(r chi.Router) {
			r.Use(auth(opts.BasicAuth, opts.Realm, opts.AuthSalt))
			r.Use(jsonp.Handler)
			r.Get("/plist-gen/{artifactID}.plist", svc.PlistGenerator)
			r.Get("/artifact-dl/{artifactID}", svc.ArtifactDownloader)
			r.Get("/artifact-icon/{name}", svc.ArtifactIcon)
			r.Get("/artifact-get-file/{artifactID}/*", svc.ArtifactGetFile)
		})
	*/

	/*
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

		// Start web server mainloop
		port := os.Getenv("PORT")
		if port == "" {
			port = "8000"
		}
		router.Run(fmt.Sprintf(":%s", port))
	*/

	return &http.Server{Handler: router}
}

/*
func getNextSongEndpoint(c *gin.Context) {
	track, err := Radio.GetNextSong()
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("# failed to get the next song: %v", err))
		return
	}

	c.String(http.StatusOK, track.Path)
}

func playlistsEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"playlists": Radio.Playlists,
	})
}

func defaultRadioEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"radio": Radio,
	})
}

func radioSkipSongEndpoint(c *gin.Context) {
	if err := Radio.SkipSong(); err != nil {
		logrus.Errorf("Failed to connect to liquidsoap: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func radioPlayTrackEndpoint(c *gin.Context) {
	var json struct {
		Hash string `form:"hash" json:"hash"`
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	track, err := Radio.GetTrackByHash(json.Hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
	}

	if err := Radio.PlayTrack(track); err != nil {
		logrus.Errorf("Failed to connect to liquidsoap: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func radioSetNextTrackEndpoint(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "not yet implemented",
	})
}

func playlistDetailEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := Radio.GetPlaylistByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"playlist": playlist,
	})
}

func trackDetailEndpoint(c *gin.Context) {
	hash := c.Param("hash")
	track, err := Radio.GetTrackByHash(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"track": track,
	})
}

func playlistUpdateEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := Radio.GetPlaylistByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	var json struct {
		SetDefault bool `form:"default" json:"default"`
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if json.SetDefault {
		Radio.DefaultPlaylist = playlist
	}

	c.JSON(http.StatusOK, gin.H{
		"playlist": playlist,
	})
}

func playlistTracksEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := Radio.GetPlaylistByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tracks": playlist.Tracks,
	})
}

type Endpoint struct {
	Source string `json:"source"`
}

func radioEndpointsEndpoint(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	mountpoints := []string{
		"mp3-192",
		"aac-192",
		"vorbis",
		"aac-128",
		"mp3-128",
	}

	endpoints := []Endpoint{}
	for _, mountpoint := range mountpoints {
		endpoint := Endpoint{
			Source: fmt.Sprintf("http://%s:4444/%s", host, mountpoint),
		}
		endpoints = append(endpoints, endpoint)
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoints": endpoints,
	})
}

func m3uPlaylistEndpoint(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	mountpoints := []string{
		"mp3-192",
		"aac-192",
		"vorbis",
		"aac-128",
		"mp3-128",
	}

	links := []string{}
	for _, mountpoint := range mountpoints {
		links = append(links, fmt.Sprintf("http://%s:4444/%s", host, mountpoint))
	}

	c.Header("Content-Type", "audio/x-mpegurl")
	c.String(http.StatusOK, strings.Join(links, "\n"))
}
*/
