package radioman

import "go.uber.org/zap"

type Opts struct {
	URL             string // self URL
	BindAddr        string
	Verbose         bool
	Logger          *zap.Logger
	LiquidsoapAddr  string
	RadioName       string
	DefaultPlaylist string
}

func NewOpts() Opts {
	return Opts{
		BindAddr:       ":8042",
		LiquidsoapAddr: "127.0.0.1:2300",
		URL:            "http://localhost:8042",
		RadioName:      "RadioMan",
		Verbose:        false,
		Logger:         nil,
	}
}

func (opts *Opts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
}
