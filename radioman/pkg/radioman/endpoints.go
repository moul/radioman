package radioman

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/jsonp"
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
	// router.Use(middleware.URLFormat)

	// public pages
	router.Route("/", func(router chi.Router) {
		router.Get("/ping", r.pingEndpoint)
		/*
			staticPrefix := "./radioman/web"
			if os.Getenv("WEBDIR") != "" {
				staticPrefix = os.Getenv("WEBDIR")
			}

			public.StaticFile("/", path.Join(staticPrefix, "static/index.html"))
			public.Static("/static", path.Join(staticPrefix, "static"))
			public.Static("/bower_components", path.Join(staticPrefix, "bower_components"))
			public.GET("/playlist.m3u", m3uPlaylistEndpoint)
		*/
	})

	// public API
	router.Route("/api", func(router chi.Router) {
		router.Use(jsonp.Handler)
		router.Get("/radios/default/endpoints", r.radioEndpointsEndpoint)
	})

	// admin pages
	router.Route("/admin", func(router chi.Router) {
		//r.Use(auth(opts.BasicAuth, opts.Realm, opts.AuthSalt))
		// FIXME: make accounts dynamic
		//accounts := gin.Accounts{"admin": "admin"}
		//admin.Use(gin.BasicAuth(accounts))
		//router.StaticFile("/admin/", path.Join(staticPrefix, "static/admin/index.html"))

		// admin API
		router.Route("/admin/api", func(router chi.Router) {
			router.Use(jsonp.Handler)
			router.Get("/playlists", r.playlistsEndpoint)
			router.Get("/playlists/:name", r.playlistDetailEndpoint)
			router.Patch("/playlists/:name", r.playlistUpdateEndpoint)
			router.Get("/playlists/:name/tracks", r.playlistTracksEndpoint)
			router.Get("/radios/default", r.defaultRadioEndpoint)
			router.Post("/radios/default/skip-song", r.radioSkipSongEndpoint)
			router.Post("/radios/default/play-track", r.radioPlayTrackEndpoint)
			router.Post("/radios/default/set-next-track", r.radioSetNextTrackEndpoint)
			router.Get("/tracks/:hash", r.trackDetailEndpoint)
		})
	})

	// liquidsoap API
	router.Route("/liq", func(router chi.Router) {
		router.Use(jsonp.Handler)
		// FIXME: auth
		router.Get("/ping", r.pingEndpoint)
		// router.Get("/get-next-song", r.getNextSongEndpoint)
	})

	docgen.PrintRoutes(router)

	return &http.Server{Handler: router}
}

func (r *Radio) pingEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

func (r *Radio) getNextSongEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
		track, err := Radio.GetNextSong()
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("# failed to get the next song: %v", err))
			return
		}

		c.String(http.StatusOK, track.Path)
	*/
}

func (r *Radio) playlistsEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
		c.JSON(http.StatusOK, gin.H{
			"playlists": Radio.Playlists,
		})
	*/
}

func (r *Radio) defaultRadioEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
		c.JSON(http.StatusOK, gin.H{
			"radio": Radio,
		})
	*/
}

func (r *Radio) radioSkipSongEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
		if err := Radio.SkipSong(); err != nil {
			r.logger.Error("failed to connect to liquidsoap", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	*/
}

func (r *Radio) radioPlayTrackEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
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
			r.logger.Error("failed to connect to liquidsoap", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	*/
}

func (r *Radio) radioSetNextTrackEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not yet implemented",
		})
	*/
}

func (r *Radio) playlistDetailEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
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
	*/
}

func (r *Radio) trackDetailEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
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
	*/
}

func (r *Radio) playlistUpdateEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
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
	*/
}

func (r *Radio) playlistTracksEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
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
	*/
}

func (r *Radio) radioEndpointsEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
		type Endpoint struct {
			Source string `json:"source"`
		}
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
	*/
}

func (r *Radio) m3uPlaylistEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
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
	*/
}
