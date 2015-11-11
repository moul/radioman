package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/kr/fs"
)

type Playlist struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	CreationDate     time.Time `json:"creation_date"`
	ModificationDate time.Time `json:"modification_date"`
	Status           string    `json:"status"`
	Stats            struct {
		Tracks int `json:"tracks"`
	} `json:"stats"`
	Tracks map[string]*Track `json:"-"`
}

type Track struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type Database struct {
	Playlists []*Playlist
}

var DB Database

func (p *Playlist) NewTrack(path string) (*Track, error) {
	if track, found := p.Tracks[path]; found {
		return track, nil
	}

	track := &Track{
		Path: path,
	}
	p.Tracks[path] = track
	p.Stats.Tracks++
	return track, nil
}

func init() {
	DB.Playlists = make([]*Playlist, 0)
	DB.NewPlaylist("manual")
	DB.NewDirectoryPlaylist("iTunes Music", "~/Music/iTunes/iTunes Media/Music/")
	DB.NewDirectoryPlaylist("iTunes Podcasts", "~/Music/iTunes/iTunes Media/Podcasts/")
	if dir, err := os.Getwd(); err == nil {
		DB.NewDirectoryPlaylist("local directory", dir)
	}
}

func (db *Database) NewPlaylist(name string) (*Playlist, error) {
	logrus.Infof("New playlist %q", name)
	playlist := &Playlist{
		Name:             name,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		Tracks:           make(map[string]*Track, 0),
		Status:           "New",
	}
	DB.Playlists = append(DB.Playlists, playlist)
	return playlist, nil
}

func (db *Database) NewDirectoryPlaylist(name string, path string) (*Playlist, error) {
	playlist, err := db.NewPlaylist(name)
	if err != nil {
		return nil, err
	}
	expandedPath, err := expandUser(path)
	if err != nil {
		return nil, err
	}
	playlist.Path = expandedPath
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
	router.GET("/api/playlists/:name/tracks", playlistTracksEndpoint)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go updatePlaylistsRoutine(&DB)

	router.Run(fmt.Sprintf(":%s", port))
}

func updatePlaylistsRoutine(db *Database) {
	for {
		for _, playlist := range db.Playlists {
			if playlist.Path == "" {
				logrus.Debugf("Playlist %q is not dynamic, skipping update", playlist.Name)
				continue
			}

			logrus.Infof("Updating playlist %q", playlist.Name)
			playlist.Status = "Updating"

			walker := fs.Walk(playlist.Path)
			for walker.Step() {
				if err := walker.Err(); err != nil {
					logrus.Warnf("walker error: %v", err)
					continue
				}
				playlist.NewTrack(walker.Path())
			}

			logrus.Infof("Playlist %q updated, %d tracks", playlist.Name, len(playlist.Tracks))
			playlist.Status = "Ready"
			playlist.ModificationDate = time.Now()
		}
		time.Sleep(5 * time.Minute)
	}
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

func playlistTracksEndpoint(c *gin.Context) {
	name := c.Param("name")
	playlist, err := DB.GetPlaylistByName(name)
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
