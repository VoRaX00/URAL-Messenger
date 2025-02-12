package handler

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log/slog"
	"messenger/internal/domain/models"
	"messenger/internal/services"
	"net/http"
	"sync"
)

type Handler struct {
	mux              *mux.Router
	wsUpg            *websocket.Upgrader
	log              *slog.Logger
	mu               sync.RWMutex
	messengerService services.IMessengerService
	clients          map[uuid.UUID]*websocket.Conn
	broadcast        chan *models.Message
}

func NewHandler(log *slog.Logger, messengerService services.IMessengerService) *Handler {
	return &Handler{
		mux: mux.NewRouter(),
		log: log,
		wsUpg: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		mu:               sync.RWMutex{},
		messengerService: messengerService,
		broadcast:        make(chan *models.Message),
	}
}

func (h *Handler) InitRoutes() {
	h.mux.HandleFunc("/ws", h.wsHandler)
	h.mux.HandleFunc("/send", h.Send)
	go h.writeToClientsBroadcast()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
