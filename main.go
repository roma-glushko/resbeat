package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"resbeat/pkg/resbeat"
	"resbeat/pkg/resbeat/telemetry"
	"syscall"
	"time"
)

// version must be set from the contents of VERSION file by go build's
// -X main.version= option in the Makefile.
var version = "unknown"

// commitSha will be the hash that the binary was built from
// and will be populated by the Makefile
var commitSha = "unknown"

func main() {
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
				Usage: "set the log format (text (default), or json)",
				Value: "text",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "set the min log level (debug, info (default), warn, error)",
				Value: "info",
			},
			&cli.DurationFlag{
				Name:  "frequency",
				Value: 3 * time.Second,
			},
			&cli.BoolFlag{
				Name:  "gpu",
				Value: false,
				Usage: "collect GPU metrics? (NVML should be installed)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			host := cCtx.String("host")
			port := cCtx.Int("port")
			logLevel := cCtx.String("log-level")
			logFormat := cCtx.String("log-format")
			frequency := cCtx.Duration("frequency")
			gpuSupport := cCtx.Bool("gpu")

			ctx := context.Background()
			ctx, logger, err := telemetry.SetupLogger(ctx, logFormat, logLevel)

			if err != nil {
				panic(err)
			}

			defer func() {
				err := logger.Sync()

				if err != nil && !errors.Is(err, syscall.ENOTTY) {
					// https://github.com/uber-go/zap/issues/991#issuecomment-962098428
					logger.Error(fmt.Sprintf("error while flushing log buffer: %v", err))
				}
			}()

			signalHandler := &resbeat.SignalHandler{}
			beatApp := resbeat.NewResBeat(ctx, gpuSupport)
			ctx = signalHandler.Handle(ctx)

			if err := beatApp.Serve(ctx, host, port, frequency); err != http.ErrServerClosed {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
