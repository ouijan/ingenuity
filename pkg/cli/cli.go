package cli

import (
	"context"
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
				configPath := cmd.Value("config").(string)
				context, err := buildContext(configPath)
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
