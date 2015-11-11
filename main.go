package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Playlist struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	CreationDate time.Time `json:"creation_date"`
	Tracks       []*Track  `json:"tracks"`
}

type Track struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type Database struct {
	Playlists []*Playlist
}

var DB Database

func init() {
	DB.Playlists = make([]*Playlist, 0)
	DB.NewPlaylist("manual")
	DB.NewDirectoryPlaylist("iTunes Music", "~/Music/iTunes/iTunes Media/Music/")
	DB.NewDirectoryPlaylist("iTunes Podcasts", "~/Music/iTunes/iTunes Media/Podcasts/")
	if dir, err := os.Getwd(); err == nil {
		DB.NewDirectoryPlaylist("local directory", dir)
	}
	//DB.NewDirectoryPlaylist("", "~/Music/iTunes/iTunes\ Media/Podcasts/")
}

func (db *Database) NewPlaylist(name string) (*Playlist, error) {
	playlist := &Playlist{
		Name:         name,
		CreationDate: time.Now(),
		Tracks:       make([]*Track, 0),
	}
	DB.Playlists = append(DB.Playlists, playlist)
	return playlist, nil
}

func (db *Database) NewDirectoryPlaylist(name string, path string) (*Playlist, error) {
	playlist, err := db.NewPlaylist(name)
	if err != nil {
		return nil, err
	}
	playlist.Path = path
	return playlist, nil
}

func (db *Database) GetPlaylistByName(name string) (*Playlist, error) {
	for _, playlist := range DB.Playlists {
		if playlist.Name == name {
			return playlist, nil
		}
	}
	return nil, fmt.Errorf("No such playlist")
}

func main() {
	router := gin.Default()

	// ping
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// static files
	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")
	router.Static("/bower_components", "./bower_components")

	router.GET("/api/playlists", playlistsEndpoint)
	router.GET("/api/playlists/:name", playlistDetailEndpoint)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(fmt.Sprintf(":%s", port))
}

func playlistsEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"playlists": DB.Playlists,
	})
}

func playlistDetailEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := DB.GetPlaylistByName(name)
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
