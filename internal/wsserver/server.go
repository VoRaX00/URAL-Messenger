package wsserver

import (
	"context"
	"errors"
	"log/slog"
	"messenger/internal/handler"
	"net/http"
)

type WSServer interface {
	MustStart()
	Start() error
	Stop(ctx context.Context)
}

type wsSrv struct {
	srv *http.Server
	log *slog.Logger
}

func NewWsServer(addr string, log *slog.Logger) WSServer {
	h := handler.NewHandler(log)
	h.InitRoutes()

	return &wsSrv{
		srv: &http.Server{
			Addr:    addr,
			Handler: h,
		},
		log: log,
	}

}

func (s *wsSrv) MustStart() {
	err := s.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *wsSrv) Start() error {
	const op = "server.Start"
	s.log.With(slog.String("op", op)).Info("starting server")
	return s.srv.ListenAndServe()
}

func (s *wsSrv) Stop(ctx context.Context) {
	const op = "server.Stop"
	s.log.With(slog.String("op", op)).Info("stopping server")

	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
}
