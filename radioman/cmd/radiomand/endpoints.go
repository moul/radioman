package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

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
