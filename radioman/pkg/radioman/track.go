package radioman

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type Track struct {
	Hash             string      `json:"hash"`
	Status           string      `json:"status"`
	Title            string      `json:"title"`
	RelPath          string      `json:"relative_path"`
	Path             string      `json:"path"`
	FileName         string      `json:"file_name"`
	FileSize         int64       `json:"file_size"`
	FileModTime      time.Time   `json:"file_modification_time"`
	CreationDate     time.Time   `json:"creation_date"`
	ModificationDate time.Time   `json:"modification_date"`
	Playlists        []*Playlist `json:"playlists"`
	Tag              struct {
		Length   time.Duration `json:"length"`
		Title    string        `json:"title"`
		Artist   string        `json:"artist"`
		Album    string        `json:"album"`
		Genre    string        `json:"genre"`
		Bitrate  int           `json:"bitrate"`
		Year     int           `json:"year"`
		Channels int           `json:"channels"`
	} `json:"tag"`
}

func NewTrack(path string) (*Track, error) {
	track := Track{
		Path:             path,
		RelPath:          path,
		Status:           "new",
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
	}
	hasher := md5.New()
	hasher.Write([]byte(path))
	track.Hash = hex.EncodeToString(hasher.Sum(nil))
	return &track, nil
}

func (t *Track) IsValid() bool {
	return t.Tag.Bitrate >= 64
}
