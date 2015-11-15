package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func getNextSongEndpoint(c *gin.Context) {
	// FIXME: shuffle playlist instead of getting a random track
	// FIXME: do not iterate over a map

	playlist := Radio.DefaultPlaylist
	track, err := playlist.GetRandomTrack()
	if err == nil {
		c.String(http.StatusOK, track.Path)
		return
	}

	for _, playlist := range Radio.Playlists {
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
		"playlists": Radio.Playlists,
	})
}

func defaultRadioEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"radio": Radio,
	})
}

func radioSkipSongEndpoint(c *gin.Context) {
	// FIXME: return json with detail
	if err := Radio.SkipSong(); err != nil {
		logrus.Errorf("Failed to connect to liquidsoap: %v", err)
		return
	}
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
