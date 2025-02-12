package services

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/internal/storage"
	"messenger/pkg/mapper"
)

type IMessengerService interface {
	Add(message domain.MessageAdd) error
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
	GetById(id uuid.UUID) (models.Message, error)
	Update(id uuid.UUID, message string) error
	Delete(id uuid.UUID) error
}

type Messenger struct {
	log        *slog.Logger
	cache      storage.MessengerCacheRepo
	repository storage.MessengerRepo
}

func NewMessenger(log *slog.Logger, cache storage.MessengerCacheRepo, repository storage.MessengerRepo) IMessengerService {
	return &Messenger{
		log:        log,
		cache:      cache,
		repository: repository,
	}
}

func (m *Messenger) Add(message domain.MessageAdd) error {
	const op = "services.messenger.Add"

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("mapping model to dto")
	dto := mapper.MessageAddToMessage(message)

	log.Info("adding message")
	err := m.repository.Add(dto)
	if err != nil {
		log.Error("error with adding message", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message added")
	return nil
}

func (m *Messenger) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	const op = "services.messenger.GetByChat"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("getting messages for chat")
	messages, err := m.repository.GetByChat(chatId)
	if err != nil {
		log.Error("error with getting messages for chat", err)
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
		log.Error("error with getting message", err)
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message received")
	return message, nil
}

func (m *Messenger) Update(id uuid.UUID, message string) error {
	const op = "services.messenger.Update"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("updating message")
	err := m.repository.Update(id, message)
	if err != nil {
		log.Error("error with updating message", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message updated")
	return nil
}

func (m *Messenger) Delete(id uuid.UUID) error {
	const op = "services.messenger.Delete"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("deleting message")
	err := m.repository.Delete(id)
	if err != nil {
		log.Error("error with deleting message", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message deleted")
	return nil
}
