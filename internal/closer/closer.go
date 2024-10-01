package closer

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
)

var globalCloser *Closer

// Add adds new function to globalCloser
func Add(stage int, f ...func() error) {
	globalCloser.Add(stage, f...)
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
	wg            *sync.WaitGroup
	mu            sync.Mutex
	once          sync.Once
	done          chan struct{}
	funcsStageOne []func() error
	funcsStageTwo []func() error
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

func (c *Closer) Add(stage int, f ...func() error) {
	c.mu.Lock()
	switch stage {
	case 1:
		c.funcsStageOne = append(c.funcsStageOne, f...)
	case 2:
		c.funcsStageTwo = append(c.funcsStageTwo, f...)
	}
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		slog.Info("running stage 1...\n")
		c.mu.Lock()
		funcs := c.funcsStageOne
		c.funcsStageOne = nil
		c.runFuncs(funcs)
		c.mu.Unlock()

		c.wg.Wait()

		slog.Info("running stage 2...\n")
		c.mu.Lock()
		funcs = c.funcsStageTwo
		c.funcsStageTwo = nil
		c.runFuncs(funcs)
		c.mu.Unlock()

		os.Exit(0)
	})
}

func (c *Closer) runFuncs(funcs []func() error) {
	errs := make(chan error, len(funcs))

	var wg sync.WaitGroup

	wg.Add(len(funcs))
	for _, f := range funcs {
		go func(f func() error) {
			defer wg.Done()

			errs <- f()
		}(f)
	}

	if len(errs) > 0 {
		for err := range errs {
			slog.Error("error returned from closer: %v\n", slog.String("error", err.Error()))
		}
	}

	wg.Wait()
}
