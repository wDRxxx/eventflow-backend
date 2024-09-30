package closer

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
)

var globalCloser *Closer

// Add adds new function to globalCloser
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait waits for executing all functions added by Add
func Wait() {
	globalCloser.Wait()
}

// CloseAll runs all functions added by Add
func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	wg    *sync.WaitGroup
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

func SetGlobalCloser(c *Closer) {
	globalCloser = c
}

func New(wg *sync.WaitGroup, sig ...os.Signal) *Closer {
	c := &Closer{
		wg:   wg,
		done: make(chan struct{}),
	}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		slog.Info("waiting for running jobs to end...")
		c.wg.Wait()
		slog.Info("done. closing all...")

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))

		c.wg.Add(len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				defer c.wg.Done()

				errs <- f()
			}(f)
		}

		if len(errs) > 0 {
			for err := range errs {
				slog.Error("error returned from closer: %v\n", slog.String("error", err.Error()))
			}
		}

		c.wg.Wait()
		os.Exit(0)
	})
}
