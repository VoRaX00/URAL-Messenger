package handler

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log/slog"
	"messenger/internal/domain"
	"net/http"
	"sync"
)

type Handler struct {
	mux       *mux.Router
	wsUpg     *websocket.Upgrader
	log       *slog.Logger
	wsClients map[*websocket.Conn]struct{}
	mu        sync.RWMutex
	broadcast chan *domain.WSMessage
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		mux: mux.NewRouter(),
		log: log,
		wsUpg: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		wsClients: make(map[*websocket.Conn]struct{}),
		mu:        sync.RWMutex{},
		broadcast: make(chan *domain.WSMessage),
	}
}

func (h *Handler) InitRoutes() {
	h.mux.HandleFunc("/ws", h.wsHandler)
	go h.writeToClientsBroadcast()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
