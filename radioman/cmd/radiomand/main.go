package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"moul.io/motd"
	"moul.io/radioman/radioman/pkg/radioman"
	"moul.io/srand"
)

func main() {
	if err := run(os.Args); err != nil {
		if err != flag.ErrHelp {
			log.Fatalf("error: %v", err)
		}
		os.Exit(1)
	}
}

func run(args []string) error {
	opts := radioman.NewOpts()

	// flags
	fs := flag.NewFlagSet("radioman", flag.ExitOnError)
	fs.BoolVar(&opts.Verbose, "v", opts.Verbose, "verbose")
	fs.StringVar(&opts.LiquidsoapAddr, "liquidsoap", opts.LiquidsoapAddr, "liquidsoap telnet TCP addr")
	fs.StringVar(&opts.BindAddr, "bind", opts.BindAddr, "HTTP bind addr")

	root := &ffcli.Command{
		ShortUsage: "radioman [flags]",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
			rand.Seed(srand.Fast())

			logger, err := zap.NewDevelopment()
			if err != nil {
				return err
			}

			opts.Logger = logger

			man, err := radioman.New(opts)
			if err != nil {
				return fmt.Errorf("radioman.New: %w", err)
			}

			err = man.Start()
			if err != nil {
				return fmt.Errorf("radioman.Start: %w", err)
			}

			return nil
		},
	}

	return root.ParseAndRun(context.Background(), args[1:])
}
