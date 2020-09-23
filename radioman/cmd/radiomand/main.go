package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/moul/radioman/radioman/pkg/radioman"
	"github.com/peterbourgon/ff/v3/ffcli"
	"moul.io/motd"
)

var opts radioman.Opts

func main() {
	if err := run(os.Args); err != nil {
		if err != flag.ErrHelp {
			log.Fatalf("error: %v", err)
		}
		os.Exit(1)
	}
}

func run(args []string) error {
	// flags
	fs := flag.NewFlagSet("radioman", flag.ExitOnError)
	fs.BoolVar(&opts.Verbose, "v", false, "verbose")

	root := &ffcli.Command{
		ShortUsage: "radioman <subcommand> [flags]",
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
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
