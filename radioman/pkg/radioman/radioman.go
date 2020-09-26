package radioman

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/oklog/run"
	"go.uber.org/zap"
	"moul.io/radioman/radioman/pkg/liquidsoap"
)

type Radio struct {
	logger    *zap.Logger
	workers   run.Group
	playlists []*Playlist
	telnet    *liquidsoap.Telnet
	config    struct {
		defaultPlaylist *Playlist
	}

	Opts    Opts
	Created time.Time
	Started time.Time
	Updated time.Time
	Stats   struct {
		Playlists int
		Tracks    int
	}
}

func New(opts Opts) (*Radio, error) {
	opts.applyDefaults()
	r := Radio{
		Opts:    opts,
		Created: time.Now(),
		Updated: time.Now(),

		logger:    opts.Logger.Named("man"),
		playlists: make([]*Playlist, 0),
	}

	// populate
	{
		err := r.stdPopulate()
		if err != nil {
			return nil, fmt.Errorf("populate error: %w", err)
		}
	}

	// web server
	{
		server := r.server()
		listener, err := net.Listen("tcp", opts.BindAddr)
		if err != nil {
			return nil, fmt.Errorf("start listener on %q: %w", opts.BindAddr, err)
		}

		r.workers.Add(func() error {
			r.logger.Info("starting HTTP server", zap.String("bind", r.Opts.BindAddr))
			return server.Serve(listener)
		}, func(err error) {
			r.logger.Info("shutting down HTTP server", zap.Error(err))

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				r.logger.Error("failed to shut down HTTP server", zap.Error(err))
			}

			_ = server.Close()
		})
	}

	// fs watcher
	{
		ctx, cancel := context.WithCancel(context.Background())
		r.workers.Add(func() error {
			return r.updatePlaylistsRoutine(ctx)
		}, func(err error) {
			cancel()
		})
	}

	// ctrl+c
	{
		ctx := context.Background()
		r.workers.Add(run.SignalHandler(ctx, syscall.SIGKILL))
	}

	// telnet
	{
		ctx := context.Background()
		r.workers.Add(func() error {
			r.logger.Info("connecting to liquidsoap telnet server", zap.String("addr", opts.LiquidsoapAddr))
			// init telnet
			r.telnet = liquidsoap.NewTelnet(opts.LiquidsoapAddr, opts.Logger)

			// send intro message on telnet
			_, err := r.telnet.Command(fmt.Sprintf(`var.set radiomand_url = %q`, opts.URL))
			if err != nil {
				return fmt.Errorf("liquidsoap error: %w", err)
			}

			// start a goroutine that send ping?
			<-ctx.Done()
			return nil
		}, func(err error) {
			// send bye message
			// close
		})
	}

	return &r, nil
}

func (r *Radio) Run() error {
	r.Started = time.Now()
	return r.workers.Run()
}
