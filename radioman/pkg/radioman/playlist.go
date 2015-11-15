package radioman

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/wtolson/go-taglib"
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

	track := &Track{
		Path:             path,
		RelPath:          relPath,
		FileName:         stat.Name(),
		FileSize:         stat.Size(),
		FileModTime:      stat.ModTime(),
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		// Mode:          stat.Mode(),
	}

	file, err := taglib.Read(path)
	if err != nil {
		logrus.Warnf("Failed to read taglib %q: %v", path, err)
	} else {
		defer file.Close()
		track.Tag.Length = file.Length()
		track.Tag.Artist = file.Artist()
		track.Tag.Title = file.Title()
		track.Tag.Album = file.Album()
		track.Tag.Genre = file.Genre()
		track.Tag.Bitrate = file.Bitrate()
		track.Tag.Year = file.Year()
		track.Tag.Channels = file.Channels()
		// fmt.Println(file.Title(), file.Artist(), file.Album(), file.Comment(), file.Genre(), file.Year(), file.Track(), file.Length(), file.Bitrate(), file.Samplerate(), file.Channels())
	}

	p.Tracks[path] = track
	p.Stats.Tracks++
	return track, nil
}

func (p *Playlist) GetTrackByPath(path string) (*Track, error) {
	if track, found := p.Tracks[path]; found {
		return track, nil
	}
	return nil, fmt.Errorf("no such track")
}

func (p *Playlist) GetRandomTrack() (*Track, error) {
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
