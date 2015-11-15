package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kr/fs"
)

type Radio struct {
	Name             string    `json:"name"`
	DefaultPlaylist  *Playlist `json:"default_playlist"`
	CreationDate     time.Time `json:"creation_date"`
	ModificationDate time.Time `json:"modification_date"`
	Stats            struct {
		Playlists int `json:"playlists"`
		Tracks    int `json:"tracks"`
	} `json:"stats"`
	Playlists []*Playlist       `json:"-"`
	Telnet    *LiquidsoapTelnet `json:"-"`
}

func (r *Radio) NewPlaylist(name string) (*Playlist, error) {
	logrus.Infof("New playlist %q", name)
	playlist := &Playlist{
		Name:             name,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		Tracks:           make(map[string]*Track, 0),
		Status:           "New",
	}
	r.Playlists = append(r.Playlists, playlist)
	r.Stats.Playlists++
	return playlist, nil
}

func (r *Radio) NewDirectoryPlaylist(name string, path string) (*Playlist, error) {
	playlist, err := r.NewPlaylist(name)
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

func (r *Radio) GetPlaylistByName(name string) (*Playlist, error) {
	for _, playlist := range r.Playlists {
		if playlist.Name == name {
			return playlist, nil
		}
	}
	return nil, fmt.Errorf("no such playlist")
}

func NewRadio(name string) *Radio {
	return &Radio{
		Name:             name,
		Playlists:        make([]*Playlist, 0),
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
	}
}

func (r *Radio) InitTelnet() error {
	liquidsoapAddr := strings.Split(strings.Replace(os.Getenv("LIQUIDSOAP_PORT_2300_TCP"), "tcp://", "", -1), ":")
	liquidsoapHost := liquidsoapAddr[0]
	liquidsoapPort, _ := strconv.Atoi(liquidsoapAddr[1])

	r.Telnet = NewLiquidsoapTelnet(liquidsoapHost, liquidsoapPort)

	if err := r.Telnet.Open(); err != nil {
		logrus.Fatalf("Failed to connect to liquidsoap")
	}
	defer r.Telnet.Close()

	radiomandHost := strings.Split(r.Telnet.Conn.LocalAddr().String(), ":")[0]
	_, err := r.Telnet.Command(fmt.Sprintf(`var.set radiomand_url = "http://%s:%d"`, radiomandHost, 8000))

	return err
}

func updatePlaylistsRoutine(r *Radio) {
	for {
		tracksSum := 0
		for _, playlist := range r.Playlists {
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
				stat := walker.Stat()

				if stat.IsDir() {
					switch stat.Name() {
					case ".git", "bower_components":
						walker.SkipDir()
					}
				} else {
					switch stat.Name() {
					case ".DS_Store":
						continue
					}

					playlist.NewLocalTrack(walker.Path())
				}
			}

			logrus.Infof("Playlist %q updated, %d tracks", playlist.Name, len(playlist.Tracks))
			playlist.Status = "Ready"
			playlist.ModificationDate = time.Now()
			tracksSum += playlist.Stats.Tracks
		}
		r.Stats.Tracks = tracksSum
		time.Sleep(5 * time.Minute)
	}
}
