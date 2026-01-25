package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var cmd = &cli.Command{
	Name:   "build_spritesheet",
	Action: runApp,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "airframes_path",
			Aliases: []string{"afp"},
			Usage:   "Path to the airframes JSON directory",
			Value:   "airframes/",
			Hidden:  true,
		},
		&cli.StringFlag{
			Name:     "inkscape_binary",
			Aliases:  []string{"inkscape"},
			Usage:    "Path to the inkscape v1+ binary",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output_png",
			Aliases:  []string{"op"},
			Usage:    "Path to the output png file",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output_json",
			Aliases:  []string{"oj"},
			Usage:    "Path to the output json file",
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
