package handler

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log/slog"
	"messenger/internal/domain/models"
	"net/http"
	"sync"
)

type Handler struct {
	mux            *mux.Router
	wsUpg          *websocket.Upgrader
	log            *slog.Logger
	mu             sync.RWMutex
	messageService MessengerService
	chatService    ChatService
	clients        map[uuid.UUID]map[uuid.UUID]*websocket.Conn
	broadcast      chan *models.Message
}

func NewHandler(log *slog.Logger, messengerService MessengerService, chatService ChatService) *Handler {
	return &Handler{
		mux: mux.NewRouter(),
		log: log,
		wsUpg: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		mu:             sync.RWMutex{},
		messageService: messengerService,
		chatService:    chatService,
		broadcast:      make(chan *models.Message),
		clients:        make(map[uuid.UUID]map[uuid.UUID]*websocket.Conn),
	}
}

func (h *Handler) InitRoutes() {
	h.mux.HandleFunc("/ws", h.wsHandler)
	h.mux.HandleFunc("/chat/add", h.AddChat).Methods("POST")
	h.mux.HandleFunc("/send", h.Send).Methods("POST")
	go h.writeToClientsBroadcast()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
