package internal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Log levels
const (
	Debug = 10
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

type Named string

func (n Named) Name() string {
	return string(n)
}

func MergeCh[T any](chs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	output := func(c <-chan T) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(chs))
	for _, c := range chs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func DebounceCh[T any](
	ch <-chan T, threshold time.Duration,
) <-chan T {
	out := make(chan T)
	go func() {
		var (
			t        = time.NewTimer(threshold)
			dispatch bool
		)
		t.Stop()

		for {
			select {
			case _, ok := <-ch:
				if !ok {
					close(out)
					return
				}

				if !dispatch {
					dispatch = true
					t.Reset(threshold)
				}

			case <-t.C:
				dispatch = false
				out <- *new(T)
			}
		}
	}()
	return out
}
