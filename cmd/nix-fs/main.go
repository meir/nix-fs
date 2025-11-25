package main

import (
	"context"
	"errors"
	"log"
	"os"

	nixfs "github.com/meir/nix-fs/internal/nix-fs"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "NixFS",
		Usage: "Manage symlinks based on NixFS state files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "state-file",
				Value: "./state.json",
				Usage: "New state file to apply",
			},
			&cli.StringFlag{
				Name:  "old-state-file",
				Value: "/etc/nix-fs.json",
				Usage: "Old state file to apply",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			stateFileLocation := cmd.String("state-file")
			if stateFileLocation == "" {
				return errors.New("state-file flag is required")
			}

			oldStateFileLocation := cmd.String("old-state-file")

			return nixfs.Run(stateFileLocation, oldStateFileLocation)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
