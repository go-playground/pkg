package appext

import (
	"context"
	. "github.com/go-playground/assert/v2"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

func TestForceExitWithNoTimeout(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	exitFn := func(code int) {
		defer wg.Done()
		Equal(t, 1, code)
	}

	c := Context().Timeout(0).ExitFn(exitFn)

	// copy of Build for testing
	var sig = make(chan os.Signal, 1)
	signal.Notify(sig, c.signals...)

	ctx, cancel := context.WithCancel(context.Background())

	go listen(sig, cancel, c.exitFn, c.timeout, c.forceExit)

	sig <- os.Interrupt
	sig <- os.Interrupt
	wg.Wait()
	Equal(t, context.Canceled, ctx.Err())
}

func TestForceExitWithTimeout(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	exitFn := func(code int) {
		defer wg.Done()
		Equal(t, 1, code)
	}

	c := Context().Timeout(time.Hour).ExitFn(exitFn)

	// copy of Build for testing
	var sig = make(chan os.Signal, 1)
	signal.Notify(sig, c.signals...)

	ctx, cancel := context.WithCancel(context.Background())

	go listen(sig, cancel, c.exitFn, c.timeout, c.forceExit)

	sig <- os.Interrupt
	sig <- os.Interrupt
	wg.Wait()
	Equal(t, context.Canceled, ctx.Err())
}

func TestTimeoutWithNoForceExit(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	exitFn := func(code int) {
		defer wg.Done()
		Equal(t, 1, code)
	}

	c := Context().Timeout(time.Millisecond * 200).ForceExit(false).ExitFn(exitFn)

	// copy of Build for testing
	var sig = make(chan os.Signal, 1)
	signal.Notify(sig, c.signals...)

	ctx, cancel := context.WithCancel(context.Background())

	go listen(sig, cancel, c.exitFn, c.timeout, c.forceExit)

	// only sending one, timeout must be reached for test to finish
	sig <- os.Interrupt
	wg.Wait()
	Equal(t, context.Canceled, ctx.Err())
}
