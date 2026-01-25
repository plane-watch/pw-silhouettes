package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var cmd = &cli.Command{
	Name:   "svg_check",
	Action: runApp,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "svg",
			Usage:    "Path to the input svg file",
			Required: true,
		},
	},
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Send()
	}
}
