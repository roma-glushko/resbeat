package resbeat

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"resbeat/pkg/resbeat/telemetry"
	"syscall"
)

type SignalHandler struct {
}

func (h *SignalHandler) Handle(ctx context.Context) context.Context {
	logger := telemetry.FromContext(ctx)
	// https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(sigs)
		cancel()
	}()

	go func() {
		select {
		case <-sigs:
			logger.Info(fmt.Sprintf("resbeat received signal (%v), terminating", "test"))
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}
