package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type Handler struct {
	mux       *mux.Router
	wsUpg     *websocket.Upgrader
	log       *slog.Logger
	wsClients map[*websocket.Conn]struct{}
	mu        sync.Mutex
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		mux:       mux.NewRouter(),
		log:       log,
		wsUpg:     &websocket.Upgrader{},
		wsClients: make(map[*websocket.Conn]struct{}),
	}
}

func (h *Handler) InitRoutes() {
	h.mux.HandleFunc("/ws", h.wsHandler)
	h.mux.HandleFunc("/test", h.testHandler)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) testHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Test is successful"))
}

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("upgrade websocket error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.log.Info(fmt.Sprintf("Client with address %s connected", conn.RemoteAddr().String()))

	h.mu.Lock()
	h.wsClients[conn] = struct{}{}
	h.mu.Unlock()

	go h.readFromClient(conn)
}

func (h *Handler) readFromClient(conn *websocket.Conn) {

}
