package handler

import (
	"encoding/json"
	"fmt"
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

//go:generate mockery --name=ChatService --output=./mocks --case=underscore
type ChatService interface {
	Add(chat domain.AddChat) (uuid.UUID, error)
	AddNewUser(chatId uuid.UUID, userId uuid.UUID) error
	RemoveUser(chatId uuid.UUID, userId uuid.UUID) error
	GetUserChats(userId uuid.UUID) ([]uuid.UUID, error)
	Update(chatId uuid.UUID) error
	Delete(chatId uuid.UUID) error
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
		_ = conn.Close()
		return
	}

	chatIds, err := h.getUserChats(userID)
	if err != nil {
		h.log.Error("get user chats error", slog.String("err", err.Error()))
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
	defer func() {
		h.disconnection(conn, userId)
	}()

	for {
		msg := new(domain.MessageAdd)
		if err := conn.ReadJSON(msg); err != nil {
			h.log.Error("Error with reading from WebSocket: ", slog.String("err", err.Error()))
			break
		}

		addedMsg, err := h.messageService.Add(*msg)
		if err != nil {
			h.log.Error("Error with adding message to Messenger: ", slog.String("err", err.Error()))
			continue
		}
		h.broadcast <- &addedMsg
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
	_ = conn.Close()
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

	msg, err := h.messageService.Add(message)
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

func (h *Handler) AddChat(w http.ResponseWriter, r *http.Request) {
	io, err := r.GetBody()
	if err != nil {
		h.log.Error("Error with reading body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var chat domain.AddChat
	decoder := json.NewDecoder(io)
	err = decoder.Decode(&chat)
	if err != nil {
		h.log.Error("Error with decoding body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := h.chatService.Add(chat)
	if err != nil {
		h.log.Error("Error with creating chat", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.mu.Lock()
	h.clients[chatId] = make(map[uuid.UUID]*websocket.Conn)
	for _, id := range chat.PersonIds {
		for _, users := range h.clients {
			if _, ok := users[id]; ok {
				h.clients[chatId][id] = users[id]
			}
		}
	}
	h.mu.Unlock()

	_, err = w.Write([]byte(fmt.Sprintf("id: %v", chatId)))
	if err != nil {
		h.log.Error("Error writing response", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Info("Success added chat to Messenger")
	w.WriteHeader(http.StatusOK)
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
	chats, err := h.chatService.GetUserChats(userID)
	if err != nil {
		h.log.Warn("Error with getting chats for user: ", slog.String("err", err.Error()))
	}
	h.log.Info("got chats for user")
	return chats, nil
}
