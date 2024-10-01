package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"syscall"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer"
)

type App struct {
	wg *sync.WaitGroup

	serviceProver *serviceProvider

	httpServer       *http.Server
	prometheusServer *http.Server

	mailer mailer.Mailer
}

func NewApp(ctx context.Context, wg *sync.WaitGroup, envPath string) (*App, error) {
	err := config.Load(envPath)
	if err != nil {
		return nil, err
	}

	app := &App{wg: wg}

	cl := closer.New(app.wg, syscall.SIGINT, syscall.SIGTERM)
	closer.SetGlobalCloser(cl)

	err = app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initDeps(ctx context.Context) error {
	a.serviceProver = newServiceProvider()

	a.initHttpServer(ctx)
	a.initMailer()

	return nil
}

func (a *App) initHttpServer(ctx context.Context) {
	s := a.serviceProver.HTTPServer(ctx, a.wg)
	a.httpServer = &http.Server{
		Addr:    a.serviceProver.HttpConfig().Address(),
		Handler: s.Handler(),
	}
}

func (a *App) initMailer() {
	m := a.serviceProver.Mailer(a.wg)
	a.mailer = m
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	// will stack if use app wg
	a.wg.Add(1)
	go func() {
		closer.Add(1, func() error {
			a.wg.Done()
			return nil
		})

		err := a.runHttpServer()
		if err != nil {
			log.Fatalf("error running http server: %v", err)
		}
	}()

	// run mailer
	a.wg.Add(1)
	go func() {
		closer.Add(1, func() error {
			a.wg.Done()
			return nil
		})

		a.mailer.ListenForMails()
	}()

	a.wg.Wait()

	return nil
}

func (a *App) runHttpServer() error {
	slog.Info("starting http server...")

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
