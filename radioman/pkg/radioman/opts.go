package radioman

import "go.uber.org/zap"

type Opts struct {
	BindAddr       string
	Verbose        bool
	Logger         *zap.Logger
	LiquidsoapAddr string
}

func NewOpts() Opts {
	return Opts{
		BindAddr:       ":8042",
		LiquidsoapAddr: "tcp://127.0.0.1:2300",
		Verbose:        false,
		Logger:         nil,
	}
}

func (opts *Opts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
}
