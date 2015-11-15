package main

import "time"

type Track struct {
	Status           string    `json:"status"`
	Title            string    `json:"title"`
	RelPath          string    `json:"relative_path"`
	Path             string    `json:"path"`
	FileName         string    `json:"file_name"`
	FileSize         int64     `json:"file_size"`
	FileModTime      time.Time `json:"file_modification_time"`
	CreationDate     time.Time `json:"creation_date"`
	ModificationDate time.Time `json:"modification_date"`
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

func (t *Track) IsValid() bool {
	return t.Tag.Bitrate >= 64
}
