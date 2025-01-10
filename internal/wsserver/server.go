package wsserver

import "net/http"

type WSServer interface {
	Start() error
	Stop() error
}

type wsSrv struct {
	mux *http.ServeMux
	srv *http.Server
}

func NewWsServer(addr string) WSServer {
	mux := http.NewServeMux()
	return &wsSrv{
		mux: mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}

}

func (s *wsSrv) Start() error {
	s.mux.HandleFunc("/test", s.testHandler)
	return s.srv.ListenAndServe()
}

func (s *wsSrv) testHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Test is successful"))
}

func (s *wsSrv) Stop() error {
	panic("implement me")
}
