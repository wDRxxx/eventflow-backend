package app

import (
	"context"
	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

type App struct {
	wg sync.WaitGroup

	serviceProver *serviceProvider

	httpServer       *http.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context, envPath string) (*App, error) {
	err := config.Load(envPath)
	if err != nil {
		return nil, err
	}

	app := &App{wg: sync.WaitGroup{}}
	err = app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initDeps(ctx context.Context) error {
	a.serviceProver = newServiceProvider()

	a.initHttpServer(ctx)

	return nil
}

func (a *App) initHttpServer(ctx context.Context) {
	s := a.serviceProver.HTTPServer(ctx)
	a.httpServer = &http.Server{
		Addr:    a.serviceProver.HttpConfig().Address(),
		Handler: s.Handler(),
	}
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		err := a.runHttpServer()
		if err != nil {
			log.Fatalf("error running http server: %v", err)
		}
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

//func (a *App) initPrometheusServer(ctx context.Context) {
//	a.prometheusServer = &http.Server{
//		Addr: a.serviceProver
//	}
//}