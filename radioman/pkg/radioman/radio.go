package radioman

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kr/fs"
	"go.uber.org/zap"
	"moul.io/godev"
	"moul.io/u"
)

func (r *Radio) updatePlaylistsRoutine(ctx context.Context) error {
	for {
		started := time.Now()
		r.logger.Debug("refreshing playlists", zap.Int("playlists", len(r.playlists)))
		tracksSum := 0
		for _, playlist := range r.playlists {
			fmt.Println(godev.PrettyJSON(playlist))
			// automatically update playlist
			if err := playlist.AutoUpdate(); err != nil {
				playlist.Status = "error"
				r.logger.Warn("failed to update playlist", zap.Error(err))
				continue
			}

			tracksSum += playlist.Stats.Tracks

			// Set default playlist if needed
			if r.config.defaultPlaylist == nil && playlist.Status == "ready" {
				r.config.defaultPlaylist = playlist
				// when getting a new playlist for the first time, a skipsong will act like the first "play"
				r.SkipSong()
			}
		}
		r.Stats.Tracks = tracksSum
		r.logger.Info("refreshed playlists", zap.Duration("duration", time.Since(started)), zap.Int("tracks", tracksSum))

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Minute):
		}
	}
}

func (r *Radio) SkipSong() error {
	_, err := r.telnet.Command("manager.skip")
	return err
}

func (r *Radio) stdPopulate() error {
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

func (r *Radio) NewPlaylist(name string) (*Playlist, error) {
	logrus.Infof("New playlist %q", name)
	playlist := &Playlist{
		Name:             name,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		Tracks:           make(map[string]*Track, 0),
		Status:           "new",
	}
	r.playlists = append(r.playlists, playlist)
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
	expandedPath, err := u.ExpandUser(path)
	if err != nil {
		return nil, err
	}
	playlist.Path = expandedPath
	return playlist, nil
}

/*

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
