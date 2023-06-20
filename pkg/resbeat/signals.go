package resbeat

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"resbeat/pkg/resbeat/telemetry"
	"sync"
	"syscall"
)

type SignalHandler struct {
	SigC chan os.Signal
	wg   *sync.WaitGroup
}

func NewSignalHandler() *SignalHandler {
	return &SignalHandler{
		SigC: make(chan os.Signal, 1),
		wg:   &sync.WaitGroup{},
	}
}

func (h *SignalHandler) Handle(ctx context.Context) context.Context {
	logger := telemetry.FromContext(ctx)
	// https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff
	ctx, cancel := context.WithCancel(ctx)

	h.wg.Add(1)

	go func() {
		signal.Notify(h.SigC, syscall.SIGINT, syscall.SIGTERM)

		defer func() {
			signal.Stop(h.SigC)
			cancel()
			h.wg.Done()
		}()

		select {
		case sig := <-h.SigC:
			logger.Info(fmt.Sprintf("resbeat received signal (%v), terminating", sig))
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	return ctx
}
