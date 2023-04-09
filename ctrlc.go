package ctrlc

import (
	"context"
	"os"
	"os/signal"
)

// Run runs a given function with a context that is canceled if one of the given signals is received.
// If no signals are listed here,
// the default is to use [os.Interrupt].
//
// If the context is canceled because a signal is caught,
// then [context.Cause] will return a [SignalError].
func Run(ctx context.Context, f func(context.Context) error, sigs ...os.Signal) error {
	if len(sigs) == 0 {
		sigs = []os.Signal{os.Interrupt}
	}

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, sigs...)
	go func() {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigch:
			cancel(SignalError{Signal: sig})
		}
	}()

	return f(ctx)
}

type SignalError struct {
	Signal os.Signal
}

func (s SignalError) Error() string {
	return "signal " + s.Signal.String()
}
