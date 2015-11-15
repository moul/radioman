package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/kr/fs"
	"github.com/moul/radioman/radioman/pkg/radioman"
)

var R *radioman.Radio

func init() {
	R = radioman.NewRadio("RadioMan")

	R.InitTelnet()

	R.NewPlaylist("manual")
	R.NewDirectoryPlaylist("iTunes Music", "~/Music/iTunes/iTunes Media/Music/")
	R.NewDirectoryPlaylist("iTunes Podcasts", "~/Music/iTunes/iTunes Media/Podcasts/")
	dir, err := os.Getwd()
	if err == nil {
		R.NewDirectoryPlaylist("local directory", dir)
	}

	for _, playlistsDir := range []string{"/playlists", path.Join(dir, "playlists")} {
		walker := fs.Walk(playlistsDir)
		for walker.Step() {
			if walker.Path() == playlistsDir {
				continue
			}
			if err := walker.Err(); err != nil {
				logrus.Warnf("walker error: %v", err)
				continue
			}

			var realpath string
			if walker.Stat().IsDir() {
				realpath = walker.Path()
				walker.SkipDir()
			} else {
				realpath, err = filepath.EvalSymlinks(walker.Path())
				if err != nil {
					logrus.Warnf("filepath.EvalSymlinks error for %q: %v", walker.Path(), err)
					continue
				}
			}

			stat, err := os.Stat(realpath)
			if err != nil {
				logrus.Warnf("os.Stat error: %v", err)
				continue
			}
			if stat.IsDir() {
				R.NewDirectoryPlaylist(fmt.Sprintf("playlist: %s", walker.Stat().Name()), realpath)
			}
		}
	}

	playlist, _ := R.GetPlaylistByName("iTunes Music")
	R.DefaultPlaylist = playlist
}

func main() {
	router := gin.Default()

	radio := R

	// ping
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// static files
	//staticPrefix := "/web"
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

	router.GET("/api/liquidsoap/getNextSong", getNextSongEndpoint)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	go radio.UpdatePlaylistsRoutine()

	router.Run(fmt.Sprintf(":%s", port))
}

func getNextSongEndpoint(c *gin.Context) {
	// FIXME: shuffle playlist instead of getting a random track
	// FIXME: do not iterate over a map

	playlist := R.DefaultPlaylist
	track, err := playlist.GetRandomTrack()
	if err == nil {
		c.String(http.StatusOK, track.Path)
		return
	}

	for _, playlist := range R.Playlists {
		track, err := playlist.GetRandomTrack()
		if err != nil {
			continue
		}
		c.String(http.StatusOK, track.Path)
		return
	}

	c.String(http.StatusNotFound, "# cannot get a random song, are your playlists empty ?")
}

func playlistsEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"playlists": R.Playlists,
	})
}

func defaultRadioEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"radio": R,
	})
}

func radioSkipSongEndpoint(c *gin.Context) {
	radio := R

	if err := radio.Telnet.Open(); err != nil {
		logrus.Errorf("Failed to connect to liquidsoap: %v", err)
		return
	}
	defer radio.Telnet.Close()

	if _, err := radio.Telnet.Command("manager.skip"); err != nil {
		logrus.Errorf("Failed to execute manager.skip: %v", err)
	}
}

func playlistDetailEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := R.GetPlaylistByName(name)
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

func playlistUpdateEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := R.GetPlaylistByName(name)
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
	}

	if json.SetDefault {
		R.DefaultPlaylist = playlist
	}

	c.JSON(http.StatusOK, gin.H{
		"playlist": playlist,
	})
}

func playlistTracksEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := R.GetPlaylistByName(name)
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
