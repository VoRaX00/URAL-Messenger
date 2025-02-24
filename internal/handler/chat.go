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

//go:generate mockery --name=ChatService --output=./mocks --case=underscore
type ChatService interface {
	Add(chat domain.AddChat) (uuid.UUID, error)
	AddNewUser(chatId uuid.UUID, userId uuid.UUID) error
	RemoveUser(chatId uuid.UUID, userId uuid.UUID) error
	GetUserChats(userId uuid.UUID) ([]uuid.UUID, error)
	Update(chat models.Chat) error
	Delete(chatId uuid.UUID) error
}

func (h *Handler) addChat(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addChat"
	log := h.log.With(
		slog.String("op", op),
	)

	io, err := r.GetBody()
	if err != nil {
		log.Error("Error with reading body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var chat domain.AddChat
	decoder := json.NewDecoder(io)
	err = decoder.Decode(&chat)
	if err != nil {
		log.Error("Error with decoding body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := h.chatService.Add(chat)
	if err != nil {
		log.Error("Error with creating chat", slog.String("err", err.Error()))
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
		log.Error("Error writing response", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("Success added chat to Messenger")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) addNewUserChat(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addNewUserChat"
	log := h.log.With(
		slog.String("op", op),
	)

	chatId, err := uuid.Parse(r.Header.Get("chatId"))
	if err != nil {
		log.Error("Error with parsing chatId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	personId, err := uuid.Parse(r.Header.Get("personId"))
	if err != nil {
		log.Error("Error with parsing personId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("adding new user")
	err = h.chatService.AddNewUser(chatId, personId)
	if err != nil {
		log.Error("Error with adding new user", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func (h *Handler) getUserChats(userID uuid.UUID) ([]uuid.UUID, error) {
	const op = "handler.getUserChats"
	log := h.log.With(
		slog.String("op", op),
	)

	log.Info("getting chats for user")
	chats, err := h.chatService.GetUserChats(userID)
	if err != nil {
		log.Error("Error with getting chats for user: ", slog.String("err", err.Error()))
	}
	log.Info("got chats for user")
	return chats, nil
}
