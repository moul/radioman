package radioman

import (
	"time"

	"moul.io/radioman/radioman/pkg/liquidsoap"
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
	Playlists []*Playlist        `json:"-"`
	Telnet    *liquidsoap.Telnet `json:"-"`
}

/*
func (r *Radio) NewPlaylist(name string) (*Playlist, error) {
	logrus.Infof("New playlist %q", name)
	playlist := &Playlist{
		Name:             name,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		Tracks:           make(map[string]*Track, 0),
		Status:           "new",
	}
	r.Playlists = append(r.Playlists, playlist)
	r.Stats.Playlists++
	return playlist, nil
}

func (r *Radio) NewDirectoryPlaylist(name string, path string) (*Playlist, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

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

func (r *Radio) GetTrackByHash(hash string) (*Track, error) {
	// FIXME: do not iterate over playlists, use a global map instead
	for _, playlist := range r.Playlists {
		if track, found := playlist.Tracks[hash]; found {
			return track, nil
		}
	}
	return nil, fmt.Errorf("no such track")
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
	if os.Getenv("LIQUIDSOAP_PORT_2300_TCP") == "" {
		return fmt.Errorf("missing LIQUIDSOAP_PORT_2300_TCP=tcp://1.2.3.4:5678")
	}
	liquidsoapAddr := strings.Split(strings.Replace(os.Getenv("LIQUIDSOAP_PORT_2300_TCP"), "tcp://", "", -1), ":")
	liquidsoapHost := liquidsoapAddr[0]
	liquidsoapPort, _ := strconv.Atoi(liquidsoapAddr[1])

	r.Telnet = liquidsoap.NewTelnet(liquidsoapHost, liquidsoapPort)

	if err := r.Telnet.Open(); err != nil {
		logrus.Fatalf("Failed to connect to liquidsoap")
	}
	defer r.Telnet.Close()

	radiomandHost := strings.Split(r.Telnet.Conn.LocalAddr().String(), ":")[0]
	_, err := r.Telnet.Command(fmt.Sprintf(`var.set radiomand_url = "http://%s:%d"`, radiomandHost, 8000))

	return err
}

func (r *Radio) SkipSong() error {
	if err := r.Telnet.Open(); err != nil {
		return err
	}
	defer r.Telnet.Close()

	if _, err := r.Telnet.Command("manager.skip"); err != nil {
		return err
	}

	return nil
}

func (r *Radio) PlayTrack(track *Track) error {
	if err := r.Telnet.Open(); err != nil {
		return err
	}
	defer r.Telnet.Close()

	if _, err := r.Telnet.Command(fmt.Sprintf("request.push %s", track.Path)); err != nil {
		return err
	}

	return nil
}

func (r *Radio) UpdatePlaylistsRoutine() {
	for {
		defaultUpdated := false
		tracksSum := 0

		for _, playlist := range r.Playlists {
			// automatically update playlist
			if err := playlist.AutoUpdate(); err != nil {
				playlist.Status = "error"
				logrus.Warnf("Failed to update playlist: %v", err)
				continue
			}

			tracksSum += playlist.Stats.Tracks

			// Set default playlist if needed
			if r.DefaultPlaylist == nil && playlist.Status == "ready" {
				r.DefaultPlaylist = playlist
				defaultUpdated = true
			}
		}

		r.Stats.Tracks = tracksSum

		// when getting a new playlist for the first time, a skipsong will act like the first "play"
		if defaultUpdated {
			r.SkipSong()
		}

		// sleep 5 minutes before next run
		time.Sleep(5 * time.Minute)
	}
}

func (r *Radio) Init() error {
	if err := r.InitTelnet(); err != nil {
		return err
	}
	return nil
}

func (r *Radio) StdPopulate() error {
	// Add a dummy manual playlist
	r.NewPlaylist("manual")

	// Add local directory
	playlistsDirs := []string{"/playlists"}
	playlistsDirs = append(playlistsDirs, path.Join(os.Getenv("HOME"), "playlists"))
	playlistsDirs = append(playlistsDirs, path.Join("/home", "playlists"))
	dir, err := os.Getwd()
	if err == nil && os.Getenv("NO_LOCAL_PLAYLISTS") != "1" {
		r.NewDirectoryPlaylist("local directory", dir)
		playlistsDirs = append(playlistsDirs, path.Join(dir, "playlists"))
	}

	// Add each folders in '/playlists' and './playlists'
	for _, playlistsDir := range playlistsDirs {
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
				r.NewDirectoryPlaylist(fmt.Sprintf("playlist: %s", walker.Stat().Name()), realpath)
			}
		}
	}

	// Add 'standard' music paths
	r.NewDirectoryPlaylist("iTunes Music", "~/Music/iTunes/iTunes Media/Music/")
	r.NewDirectoryPlaylist("iTunes Podcasts", "~/Music/iTunes/iTunes Media/Podcasts/")
	r.NewDirectoryPlaylist("iTunes Music", "/home/Music/iTunes/iTunes Media/Music/")
	r.NewDirectoryPlaylist("iTunes Podcasts", "/home/Music/iTunes/iTunes Media/Podcasts/")

	return nil
}

func (r *Radio) GetNextSong() (*Track, error) {
	// FIXME: shuffle playlist instead of getting a random track
	// FIXME: do not iterate over a map

	playlist := r.DefaultPlaylist
	track, err := playlist.GetRandomTrack()
	if err == nil {
		return track, nil
	}

	for _, playlist := range r.Playlists {
		track, err := playlist.GetRandomTrack()
		if err != nil {
			continue
		}
		return track, nil
	}

	return nil, fmt.Errorf("no such next song, are your playlists empty ?")
}
*/
