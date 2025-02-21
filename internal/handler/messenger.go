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

//go:generate mockery --name=MessengerService --output=./mocks --case=underscore
type MessengerService interface {
	Add(message domain.MessageAdd) (models.Message, error)
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
	GetById(id uuid.UUID) (models.Message, error)
	Update(message domain.MessageUpdate) error
	Delete(id uuid.UUID) error
}

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("upgrade websocket error", "error", err)
		return
	}
	h.log.Info(fmt.Sprintf("Client with address %s connected", conn.RemoteAddr().String()))
	go h.Conn(conn)
}

func (h *Handler) Conn(conn *websocket.Conn) {

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

		h.mu.RLock()
		_, ok := h.clients[addedMsg.Chat.Id]
		if !ok {
			h.clients[addedMsg.Chat.Id] = make(map[*websocket.Conn]struct{})
		}
		h.clients[addedMsg.Chat.Id][conn] = struct{}{}
		h.mu.RUnlock()

		h.broadcast <- &addedMsg
		go h.writeToClientsBroadcast()
	}

	for key := range h.clients {
		if _, ok := h.clients[key][conn]; ok {
			delete(h.clients[key], conn)
		}
	}

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
		clients := h.clients[msg.Chat.Id]
		for client, _ := range clients {
			if err := client.WriteJSON(msg); err != nil {
				h.log.Warn("Error with adding message to Messenger: ", slog.String("err", err.Error()))
			}
		}
		h.mu.RUnlock()
	}
}
