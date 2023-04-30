package main

import (
	"fmt"

	"os"

	"github.com/urfave/cli/v2"
)

var version = "git"

func main() {
	app := &cli.App{
		Name:    "go-matrix-webhook",
		Version: version,
		Usage:   "Matrix Webhook",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listen-addr",
				Usage:   "listen addr",
				EnvVars: []string{"LISTEN_ADDR"},
				Value:   "0.0.0.0",
			},
			&cli.IntFlag{
				Name:    "listen-port",
				Usage:   "listen port",
				EnvVars: []string{"LISTEN_PORT"},
				Value:   8080,
			},
			&cli.StringFlag{
				Name:    "listen-path",
				Usage:   "listen path",
				EnvVars: []string{"LISTEN_PATH"},
				Value:   "/",
			},
			&cli.StringFlag{
				Name:     "secret-header",
				Usage:    "secret header",
				EnvVars:  []string{"SECRET_HEADER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "matrix-access-token",
				Usage:    "matrix access token",
				EnvVars:  []string{"MATRIX_ACCESS_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "matrix-id",
				Usage:    "matrix id",
				EnvVars:  []string{"MATRIX_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "matrix-url",
				Usage:    "matrix url",
				EnvVars:  []string{"MATRIX_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "logger level (debug, info, warn, error, fatal, panic, no)",
				EnvVars: []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-type",
				Value:   "pretty",
				Usage:   "logger type (pretty, json)",
				EnvVars: []string{"LOG_TYPE"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\n" + err.Error())
	}
}
