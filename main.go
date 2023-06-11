package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"resbeat/pkg/resbeat"
	"resbeat/pkg/resbeat/telemetry"
	"time"
)

// version must be set from the contents of VERSION file by go build's
// -X main.version= option in the Makefile.
var version = "unknown"

// commitSha will be the hash that the binary was built from
// and will be populated by the Makefile
var commitSha = "unknown"

var ctx context.Context
var signalHandler *resbeat.SignalHandler
var beatApp *resbeat.ResBeat

func init() {
	ctx = context.Background()
	signalHandler = &resbeat.SignalHandler{}
	beatApp = resbeat.NewResBeat()
}

func main() {
	logger, err := telemetry.SetupLogger(ctx)

	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:      "resbeat",
		Usage:     "🔊 broadcast container resource utilization via HTTP polling or websocket",
		Copyright: "Roman Hlushko, 2023",
		Version:   resbeat.GetVersion(version, commitSha),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "127.0.0.1",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 8000,
			},
			&cli.StringFlag{
				Name:  "log-format",
				Usage: "set the log format ('text' (default), or 'json')",
				Value: "text",
			},
			&cli.DurationFlag{
				Name:  "frequency",
				Value: 5 * time.Second,
			},
		},
		Action: func(cCtx *cli.Context) error {
			host := cCtx.String("host")
			port := cCtx.Int("port")
			frequency := cCtx.Duration("frequency")

			cancelCtx := signalHandler.Handle(ctx)

			if err := beatApp.Serve(cancelCtx, host, port, frequency); err != http.ErrServerClosed {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}
