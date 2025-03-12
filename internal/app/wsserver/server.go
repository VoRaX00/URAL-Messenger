package wsserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type WSServer interface {
	MustStart()
	Start() error
	MustStop(ctx context.Context)
	Stop(ctx context.Context) error
}

type wsSrv struct {
	srv *http.Server
	log *slog.Logger
}

func New(log *slog.Logger, h http.Handler, config Config) WSServer {
	return &wsSrv{
		srv: &http.Server{
			Addr:         config.Addr,
			Handler:      h,
			ReadTimeout:  config.Timeout,
			WriteTimeout: config.Timeout,
		},
		log: log,
	}
}

func (s *wsSrv) MustStart() {
	err := s.Start()
	if err != nil {
		panic(err)
	}
}

func (s *wsSrv) Start() error {
	const op = "server.Start"
	s.log.With(slog.String("op", op)).Info("starting server")

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.srv.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("server failed")
			return err
		}
	case <-time.After(time.Second):
		s.log.Info("server started successfully")
	}
	return nil
}

func (s *wsSrv) MustStop(ctx context.Context) {
	const op = "server.Stop"
	s.log.With(slog.String("op", op)).Info("stopping server")
	err := s.Stop(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *wsSrv) Stop(ctx context.Context) error {
	s.log.Info("attempting graceful shutdown...")

	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Error("failed to shutdown server", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("server successfully shutdown")
	return nil
}
