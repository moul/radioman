package radioman

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kr/fs"
	taglib "github.com/wtolson/go-taglib"
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

func (p *Playlist) NewLocalTrack(path string) (*Track, error) {
	if track, err := p.GetTrackByPath(path); err == nil {
		return track, nil
	}

	relPath := path
	if strings.Index(path, p.Path) == 0 {
		relPath = path[len(p.Path):]
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	track, err := NewTrack(path)
	if err != nil {
		return nil, err
	}
	track.RelPath = relPath
	track.FileName = stat.Name()
	track.FileSize = stat.Size()
	track.FileModTime = stat.ModTime()

	file, err := taglib.Read(path)
	if err != nil {
		logrus.Warnf("Failed to read taglib %q: %v", path, err)
		track.Status = "error"
		track.Title = track.FileName
	} else {
		defer file.Close()
		track.Tag.Length = file.Length() / time.Second
		track.Tag.Artist = file.Artist()
		track.Tag.Title = file.Title()
		track.Tag.Album = file.Album()
		track.Tag.Genre = file.Genre()
		track.Tag.Bitrate = file.Bitrate()
		track.Tag.Year = file.Year()
		track.Tag.Channels = file.Channels()
		// FIXME: do not prepend the artist if it is already present in the title
		track.Title = fmt.Sprintf("%s - %s", track.Tag.Artist, track.Tag.Title)
		track.Status = "ready"
		// fmt.Println(file.Title(), file.Artist(), file.Album(), file.Comment(), file.Genre(), file.Year(), file.Track(), file.Length(), file.Bitrate(), file.Samplerate(), file.Channels())
	}

	p.Tracks[track.Hash] = track
	p.Stats.Tracks++
	return track, nil
}

func (p *Playlist) GetTrackByPath(path string) (*Track, error) {
	// FIXME: use a dedicated map
	for _, track := range p.Tracks {
		if track.Path == path {
			return track, nil
		}
	}
	return nil, fmt.Errorf("no such track")
}

func (p *Playlist) GetRandomTrack() (*Track, error) {
	if p == nil {
		return nil, fmt.Errorf("playlist is nil")
	}
	if p.Status != "ready" {
		return nil, fmt.Errorf("playlist is not ready")
	}

	validFiles := 0
	for _, track := range p.Tracks {
		if track.IsValid() {
			validFiles++
		}
	}

	if validFiles == 0 {
		return nil, fmt.Errorf("there is no available track")
	}

	i := rand.Intn(validFiles)
	for _, track := range p.Tracks {
		if !track.IsValid() {
			continue
		}
		if i <= 0 {
			return track, nil
		}
		i--
	}

	return nil, fmt.Errorf("cannot get a random track")
}

func (p *Playlist) AutoUpdate() error {
	if p.Path == "" {
		logrus.Debugf("Playlist %q is not dynamic, skipping update", p.Name)
		return nil
	}

	// if we are here, the playlist is based on local file system
	logrus.Infof("Updating playlist %q", p.Name)

	p.Status = "updating"

	walker := fs.Walk(p.Path)

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

			p.NewLocalTrack(walker.Path())
		}
	}

	logrus.Infof("Playlist %q updated, %d tracks", p.Name, len(p.Tracks))
	if p.Stats.Tracks > 0 {
		p.Status = "ready"
	} else {
		p.Status = "empty"
	}
	p.ModificationDate = time.Now()

	return nil
}
