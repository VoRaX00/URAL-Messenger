package wsserver

import (
	"log/slog"
	"messenger/internal/handler"
	"net/http"
)

type WSServer interface {
	Start() error
	Stop() error
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

func (s *wsSrv) Start() error {
	return s.srv.ListenAndServe()
}

func (s *wsSrv) Stop() error {
	panic("implement me")
}
