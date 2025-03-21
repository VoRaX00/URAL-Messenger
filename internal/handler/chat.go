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
	"strconv"
)

//go:generate mockery --name=ChatService --output=./mocks --case=underscore
type ChatService interface {
	Add(chat domain.AddChat) (uuid.UUID, error)
	AddNewUser(chatId uuid.UUID, userId uuid.UUID) error
	RemoveUser(chatId uuid.UUID, userId uuid.UUID) error
	GetInfoUserChats(userId uuid.UUID, page, count uint) ([]domain.GetChat, error)
	GetUserChats(userId uuid.UUID) ([]uuid.UUID, error)
	GetUsers(chatId uuid.UUID) ([]uuid.UUID, error)
	GetUserInfo(id uuid.UUID) (domain.UserInfo, error)
	Update(chat models.Chat) error
	Delete(chatId, userId uuid.UUID) error
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

	chatId, err := uuid.Parse(r.URL.Query().Get("chatId"))
	if err != nil {
		log.Error("Error with parsing chatId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	personId, err := uuid.Parse(r.URL.Query().Get("personId"))
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

func (h *Handler) getInfoUserChats(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getInfoUserChats"
	log := h.log.With(
		slog.String("op", op),
	)

	userId, err := uuid.Parse(r.URL.Query().Get("userId"))
	if err != nil {
		log.Error("Error with parsing userId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		log.Error("Error with parsing page")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	countChats, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || countChats < 1 {
		log.Error("Error with parsing count")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("getting info about chats")
	chats, err := h.chatService.GetInfoUserChats(userId, uint(page), uint(countChats))
	if err != nil {
		log.Error("Error with getting info", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("got info about chats")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(chats); err != nil {
		log.Error("Error writing response", slog.String("err", err.Error()))
	}
}

func (h *Handler) getPersons(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getPersons"
	log := h.log.With(
		slog.String("op", op),
	)

	chatId, err := uuid.Parse(r.URL.Query().Get("chatId"))
	if err != nil {
		log.Error("Error with parsing chatId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("getting chats for user")
	ids, err := h.chatService.GetUsers(chatId)
	if err != nil {
		log.Error("Error with getting persons", slog.String("err", err.Error()))
	}
	log.Info("got chats for user")

	log.Info("getting info about user")
	users := make([]domain.UserInfo, len(ids))
	for i, id := range ids {
		users[i], err = h.chatService.GetUserInfo(id)
		if err != nil {
			log.Error("Error with getting info", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	log.Info("got info about user")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) removeUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.removeUser"
	log := h.log.With(
		slog.String("op", op),
	)

	chatId, err := uuid.Parse(r.URL.Query().Get("chatId"))
	if err != nil {
		log.Error("Error with parsing chatId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(r.URL.Query().Get("userId"))
	if err != nil {
		log.Error("Error with parsing userId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("removing user")
	err = h.chatService.RemoveUser(chatId, userId)
	if err != nil {
		log.Error("Error with removing user", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	const op = "handler.update"
	log := h.log.With(
		slog.String("op", op),
	)

	var chat models.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		log.Error("Error with parsing chat", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("updating chat", slog.String("chatId", chat.Id.String()))
	err = h.chatService.Update(chat)
	if err != nil {
		log.Error("Error with updating chat", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("successfully updated chat", slog.String("chatId", chat.Id.String()))

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	const op = "handler.delete"
	log := h.log.With(
		slog.String("op", op),
	)

	chatId, err := uuid.Parse(r.URL.Query().Get("chatId"))
	if err != nil {
		log.Error("Error with parsing chatId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(r.URL.Query().Get("userId"))
	if err != nil {
		log.Error("Error with parsing userId", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("deleting chat", slog.String("chatId", chatId.String()))
	err = h.chatService.Delete(chatId, userId)
	if err != nil {
		log.Error("Error with deleting chat", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("successfully deleted chat", slog.String("chatId", chatId.String()))

	w.WriteHeader(http.StatusOK)
}
