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

	logger.Info("resbeat is starting")

	app := &cli.App{
		Name:  "resbeat",
		Usage: "🔊 broadcast container resource utilization via websocket",
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
				Name:  "logformat",
				Value: "ecs",
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
				logger.Fatal(err.Error())
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}
