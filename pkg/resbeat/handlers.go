package resbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"net/http"
	"resbeat/pkg/resbeat/readers"
	"resbeat/pkg/resbeat/telemetry"
	"time"
)

type ResBeat struct {
	ctx     context.Context
	melody  *melody.Melody
	sig     *SignalHandler
	monitor *Monitor
	encoder *json.Encoder
}

func NewResBeat() *ResBeat {
	return &ResBeat{
		melody:  melody.New(),
		sig:     &SignalHandler{},
		monitor: NewMonitor(readers.DummyStatsReader{}), // TODO: use strategy to select real readers
		encoder: &json.Encoder{},
	}
}

type ShutdownFunc func() error

func (b *ResBeat) Serve(ctx context.Context, host string, port int, frequency time.Duration) error {
	logger := telemetry.FromContext(b.ctx)
	beatC := b.monitor.Run(ctx, frequency)

	srv := http.Server{Addr: fmt.Sprintf("%s:%d", host, port)}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		b.melody.HandleRequest(w, r)
	})

	b.melody.HandleConnect(func(s *melody.Session) {
		logger.Info("client connected")
	})

	b.melody.HandleDisconnect(func(s *melody.Session) {
		logger.Info("client disconnected")
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-beatC:
				usage := b.monitor.Usage()
				usageEncoded, err := json.Marshal(usage)

				if err != nil {
					logger.Warn(fmt.Sprintf("failed to encode utilization: %s", err))
					continue
				}

				logger.Debug(fmt.Sprintf("resource utilization updated: %v", string(usageEncoded)))

				err = b.melody.Broadcast(usageEncoded)

				if err != nil {
					logger.Warn(fmt.Sprintf("failed to broacast utilization: %s", err))
				}
			}
		}
	}()

	logger.Info(fmt.Sprintf("utilization beat is available at ws://%v:%d/ws", host, port))

	go func() {
		<-ctx.Done()

		logger.Info("server is shutting down")

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("server failed to shutdown: %v", err))
		}
	}()

	return srv.ListenAndServe()
}
