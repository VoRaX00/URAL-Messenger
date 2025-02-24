package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"net/http"
)

//go:generate mockery --name=MessageService --output=./mocks --case=underscore
type MessageService interface {
	Add(message domain.MessageAdd) (models.Message, error)
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
	GetById(id uuid.UUID) (models.Message, error)
	Update(message domain.MessageUpdate) error
	Delete(id uuid.UUID) error
}

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.wsHandler"
	log := h.log.With(
		slog.String("op", op),
	)

	conn, err := h.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade websocket error", "error", err)
		return
	}

	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		log.Error("parse user id error", slog.String("err", err.Error()))
		_ = conn.Close()
		return
	}

	chatIds, err := h.getUserChats(userID)
	if err != nil {
		log.Error("get user chats error", slog.String("err", err.Error()))
		_ = conn.Close()
		return
	}

	h.mu.Lock()
	for _, chatId := range chatIds {
		if h.clients[chatId] == nil {
			h.clients[chatId] = make(map[uuid.UUID]*websocket.Conn)
		}
		h.clients[chatId][userID] = conn
	}
	h.mu.Unlock()

	go h.Conn(conn, userID)
}

func (h *Handler) Conn(conn *websocket.Conn, userId uuid.UUID) {
	const op = "handler.Conn"
	log := h.log.With(
		slog.String("op", op),
	)

	defer func() {
		log.Info("User disconnection: ", slog.String("userId", userId.String()))
		h.disconnection(conn, userId)
	}()

	for {
		msg := new(domain.MessageAdd)
		if err := conn.ReadJSON(msg); err != nil {
			log.Error("Error with reading from WebSocket: ", slog.String("err", err.Error()))
			break
		}

		addedMsg, err := h.messageService.Add(*msg)
		if err != nil {
			log.Error("Error with adding message to Messenger: ", slog.String("err", err.Error()))
			continue
		}
		h.broadcast <- &addedMsg
	}
}

func (h *Handler) disconnection(conn *websocket.Conn, userId uuid.UUID) {
	const op = "handler.disconnection"
	log := h.log.With(
		slog.String("op", op),
	)

	h.mu.Lock()
	defer h.mu.Unlock()

	for chatId, subscribers := range h.clients {
		if _, ok := subscribers[userId]; ok {
			delete(h.clients, chatId)
			if len(h.clients) == 0 {
				delete(h.clients, chatId)
			}
		}
	}
	_ = conn.Close()
	log.Info("close websocket connection")
}

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	const op = "handler.send"
	log := h.log.With(
		slog.String("op", op),
	)

	io, err := r.GetBody()
	if err != nil {
		log.Error("Error with reading body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message domain.MessageAdd
	decoder := json.NewDecoder(io)
	err = decoder.Decode(&message)
	if err != nil {
		log.Error("Error with decoding body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err := h.messageService.Add(message)
	if err != nil {
		log.Error("Error with adding message to Messenger: ", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.broadcast <- &msg

	log.Info("Success added message to Messenger")
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) writeToClientsBroadcast() {
	const op = "handler.writeToClientsBroadcast"
	log := h.log.With(
		slog.String("op", op),
	)

	for msg := range h.broadcast {
		h.mu.RLock()
		sub, ok := h.clients[msg.Chat.Id]
		if ok {
			for _, client := range sub {
				if err := client.WriteJSON(msg); err != nil {
					log.Warn("Error with adding message to Messenger: ", slog.String("err", err.Error()))
				}
			}
		}
		h.mu.RUnlock()
	}
}
