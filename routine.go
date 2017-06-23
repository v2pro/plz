package plz

import (
	"context"
	"github.com/v2pro/plz/routine"
)

func Go(oneOff func()) error {
	return routine.Go(oneOff)
}

func GoLongRunning(longRunning func(ctx context.Context)) (context.CancelFunc, error) {
	return routine.GoLongRunning(longRunning)
}
