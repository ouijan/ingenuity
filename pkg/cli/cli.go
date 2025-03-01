package cli

import (
	"context"
	"errors"
	"os"

	"github.com/urfave/cli/v3"
)

var app = &cli.Command{
	Name:                  "Ingenuity CLI",
	EnableShellCompletion: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "ingenuity.json",
		},
	},
	Commands: []*cli.Command{
		{
			Name:    "tiled",
			Aliases: []string{"t"},
			Usage:   "generate Commands, Custom Enums & Custom Properties from Golang source files",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				value := cmd.Value("config")
				if value == nil {
					return errors.New("config flag is required")
				}
				path, ok := value.(string)
				if !ok {
					return errors.New("config flag must be a string")
				}
				context, err := buildContext(path)
				if err != nil {
					return err
				}
				return generateTiled(context)
			},
		},
	},
}

func Run() error {
	return app.Run(context.Background(), os.Args)
}
