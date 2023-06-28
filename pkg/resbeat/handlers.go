package resbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"go.uber.org/zap"
	"net/http"
	"resbeat/pkg/resbeat/readers/system"
	"resbeat/pkg/resbeat/telemetry"
	"time"
)

type HTTPError struct {
	Message string `json:"message"`
}

type ResBeat struct {
	melody  *melody.Melody
	sig     *SignalHandler
	monitor *Monitor
	encoder *json.Encoder
}

func NewResBeat(ctx context.Context) *ResBeat {
	logger := telemetry.FromContext(ctx)
	systemReader, err := system.NewSystemReader(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("could not init a system stat reader: %v", err))
	}

	return &ResBeat{
		melody:  melody.New(),
		sig:     NewSignalHandler(),
		monitor: NewMonitor(systemReader),
		encoder: &json.Encoder{},
	}
}

type ShutdownFunc func() error

func (b *ResBeat) Serve(ctx context.Context, host string, port int, frequency time.Duration) error {
	logger := telemetry.FromContext(ctx)
	beatC := b.monitor.Run(ctx, frequency)

	logger.Info("resbeat is starting")

	srv := http.Server{Addr: fmt.Sprintf("%s:%d", host, port)}

	// websocket API
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		err := b.melody.HandleRequest(w, r)

		if err != nil {
			logger.Error(err.Error())
		}
	})

	// HTTP Polling API
	http.HandleFunc("/usage/", func(w http.ResponseWriter, r *http.Request) {
		jsonEncoder := json.NewEncoder(w)
		usage := b.monitor.Usage()

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		err := jsonEncoder.Encode(usage)

		if err != nil {
			errMsg := fmt.Sprintf("error building the response: %v", err)
			logger.Error(errMsg)

			w.WriteHeader(http.StatusInternalServerError)
			err = jsonEncoder.Encode(HTTPError{Message: errMsg})

			if err != nil {
				logger.Error(err.Error())
			}

			return
		}
	})

	b.melody.HandleConnect(func(s *melody.Session) {
		logger.Info(
			"websocket client connected",
			zap.String("remoteAddr", s.RemoteAddr().String()),
		)
	})

	b.melody.HandleDisconnect(func(s *melody.Session) {
		logger.Info(
			"websocket client disconnected",
			zap.String("remoteAddr", s.RemoteAddr().String()),
		)
	})

	b.melody.HandleError(func(s *melody.Session, err error) {

		logger.Warn(
			fmt.Sprintf("websocket client error: %v", err),
			zap.String("remoteAddr", s.RemoteAddr().String()),
		)
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

	logger.Info(fmt.Sprintf(
		"resbeat is up and running "+
			"\n • Websocket API: ws://%v:%d/ws/ "+
			"\n • HTTP Polling API: http://%v:%d/usage/",
		host, port,
		host, port,
	))

	go func() {
		<-ctx.Done()

		logger.Info("resbeat is shutting down")

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("server failed to shutdown: %v", err))
		}
	}()

	return srv.ListenAndServe()
}
