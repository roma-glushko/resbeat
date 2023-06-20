package resbeat

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
)

func TestSignalHandler_CatchSignals(t *testing.T) {
	tests := map[string]struct {
		signal os.Signal
	}{
		"handle SIGINT": {
			signal: syscall.SIGINT,
		},
		"handle SIGTERM": {
			signal: syscall.SIGTERM,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			sigHandler := NewSignalHandler()

			ctx = sigHandler.Handle(ctx)

			sigHandler.SigC <- test.signal
			sigHandler.wg.Wait()

			assert.ErrorIs(t, ctx.Err(), context.Canceled)
		})
	}
}

func TestSignalHandler_ExitOnCtxCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigHandler := NewSignalHandler()

	ctx = sigHandler.Handle(ctx)
	cancel()

	sigHandler.wg.Wait()
}
