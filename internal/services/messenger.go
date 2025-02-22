package services

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/pkg/mapper"
)

type MessengerCacheRepo interface {
	Add(message models.Message) error
	Delete(message models.Message) error
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
}

type MessengerRepo interface {
	Add(message models.Message) (models.Message, error)
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
	GetById(id uuid.UUID) (models.Message, error)
	GetUserChats(userId uuid.UUID) ([]uuid.UUID, error)
	Update(message domain.MessageUpdate) error
	Delete(id uuid.UUID) error
}

type Messenger struct {
	log        *slog.Logger
	cache      MessengerCacheRepo
	repository MessengerRepo
}

func NewMessenger(log *slog.Logger, cache MessengerCacheRepo, repository MessengerRepo) *Messenger {
	return &Messenger{
		log:        log,
		cache:      cache,
		repository: repository,
	}
}

func (m *Messenger) Add(message domain.MessageAdd) (models.Message, error) {
	const op = "services.messenger.Add"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("mapping model to dto")
	dto := mapper.MessageAddToMessage(message)

	log.Info("adding message to relation db")
	msg, err := m.repository.Add(dto)
	if err != nil {
		log.Warn("error with adding message to relation db", err)
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message added in relation db")

	log.Info("adding message to cache")
	err = m.cache.Add(dto)
	if err != nil {
		log.Warn("error with adding message to cache", err)
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message added in cache")

	return msg, nil
}

func (m *Messenger) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	const op = "services.messenger.GetByChat"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("getting messages for chat")
	messages, err := m.repository.GetByChat(chatId)
	if err != nil {
		log.Warn("error with getting messages for chat", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("messages received")
	return messages, nil
}

func (m *Messenger) GetById(id uuid.UUID) (models.Message, error) {
	const op = "services.messenger.GetById"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("getting message")
	message, err := m.repository.GetById(id)
	if err != nil {
		log.Warn("error with getting message", err)
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message received")
	return message, nil
}

func (m *Messenger) Update(message domain.MessageUpdate) error {
	const op = "services.messenger.Update"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("updating message")
	err := m.repository.Update(message)
	if err != nil {
		log.Warn("error with updating message", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message updated")
	return nil
}

func (m *Messenger) GetUserChats(userId uuid.UUID) ([]uuid.UUID, error) {
	const op = "services.messenger.GetUserChats"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("getting chats for user")
	chats, err := m.repository.GetUserChats(userId)
	if err != nil {
		log.Error("error with getting chats for user:", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("chats received")
	return chats, nil
}

func (m *Messenger) Delete(id uuid.UUID) error {
	const op = "services.messenger.Delete"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("deleting message")
	err := m.repository.Delete(id)
	if err != nil {
		log.Warn("error with deleting message", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message deleted")
	return nil
}
