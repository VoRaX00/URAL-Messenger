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

//go:generate mockery --name=MessengerService --output=./mocks --case=underscore
type MessengerService interface {
	Add(message domain.MessageAdd) (models.Message, error)
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
	GetById(id uuid.UUID) (models.Message, error)
	Update(message domain.MessageUpdate) error
	Delete(id uuid.UUID) error
	GetUserChats(userId uuid.UUID) ([]uuid.UUID, error)
}

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("upgrade websocket error", "error", err)
		return
	}

	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		h.log.Error("parse user id error", slog.String("err", err.Error()))
		conn.Close()
		return
	}

	chatIds, err := h.getUserChats(userID)
	if err != nil {
		h.log.Error("get user chats error", slog.String("err", err.Error()))
		conn.Close()
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
	defer func() {
		h.disconnection(conn, userId)
	}()

	for {
		msg := new(domain.MessageAdd)
		if err := conn.ReadJSON(msg); err != nil {
			h.log.Warn("Error with reading from WebSocket: ", slog.String("err", err.Error()))
			break
		}

		addedMsg, err := h.messengerService.Add(*msg)
		if err != nil {
			h.log.Warn("Error with adding message to Messenger: ", slog.String("err", err.Error()))
		}

		h.broadcast <- &addedMsg
		go h.writeToClientsBroadcast()
	}
}

func (h *Handler) disconnection(conn *websocket.Conn, userId uuid.UUID) {
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
	conn.Close()
	h.log.Info("close websocket connection")
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	io, err := r.GetBody()
	if err != nil {
		h.log.Warn("Error with reading body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message domain.MessageAdd
	decoder := json.NewDecoder(io)
	err = decoder.Decode(&message)
	if err != nil {
		h.log.Warn("Error with decoding body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err := h.messengerService.Add(message)
	if err != nil {
		h.log.Warn("Error with adding message to Messenger: ", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.broadcast <- &msg

	h.log.Info("Success added message to Messenger")
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) writeToClientsBroadcast() {
	for msg := range h.broadcast {
		h.mu.RLock()
		sub, ok := h.clients[msg.Chat.Id]
		if ok {
			for _, client := range sub {
				if err := client.WriteJSON(msg); err != nil {
					h.log.Warn("Error with adding message to Messenger: ", slog.String("err", err.Error()))
				}
			}
		}
		h.mu.RUnlock()
	}
}

func (h *Handler) getUserChats(userID uuid.UUID) ([]uuid.UUID, error) {
	h.log.Info("getting chats for user")
	chats, err := h.messengerService.GetUserChats(userID)
	if err != nil {
		h.log.Warn("Error with getting chats for user: ", slog.String("err", err.Error()))
	}
	h.log.Info("got chats for user")
	return chats, nil
}
