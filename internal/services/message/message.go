package message

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/pkg/mapper"
)

//go:generate mockery --name=CacheRepository --output=./mocks --case=underscore
type CacheRepository interface {
	Add(message models.Message) error
	Delete(message models.Message) error
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
}

//go:generate mockery --name=Repository --output=./mocks --case=underscore
type Repository interface {
	Add(message models.Message) (models.Message, error)
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
	GetById(id uuid.UUID) (models.Message, error)
	Update(message models.Message) error
	Delete(id uuid.UUID) error
}

type Service struct {
	log        *slog.Logger
	cache      CacheRepository
	repository Repository
}

func NewMessageService(log *slog.Logger, cache CacheRepository, repository Repository) *Service {
	return &Service{
		log:        log,
		cache:      cache,
		repository: repository,
	}
}

func (m *Service) Add(message domain.MessageAdd) (models.Message, error) {
	const op = "services.messenger.Add"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("mapping model to dto")
	dto := mapper.MessageAddToMessage(message)

	log.Info("adding message to relation db")
	msg, err := m.repository.Add(dto)
	if err != nil {
		log.Error("error with adding message to relation db", slog.String("err", err.Error()))
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message added in relation db")

	log.Info("adding message to cache")
	err = m.cache.Add(dto)
	if err != nil {
		log.Error("error with adding message to cache", slog.String("err", err.Error()))
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message added in cache")

	return msg, nil
}

func (m *Service) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	const op = "services.messenger.GetByChat"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("getting messages for chat")
	messages, err := m.repository.GetByChat(chatId)
	if err != nil {
		log.Error("error with getting messages for chat", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("messages received")
	return messages, nil
}

func (m *Service) GetById(id uuid.UUID) (models.Message, error) {
	const op = "services.messenger.GetById"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("getting message")
	message, err := m.repository.GetById(id)
	if err != nil {
		log.Error("error with getting message", slog.String("err", err.Error()))
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message received")
	return message, nil
}

func (m *Service) Update(message domain.MessageUpdate) error {
	const op = "services.messenger.Update"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("mapping model to dto")
	dto := mapper.MessageUpdateToMessage(message)

	log.Info("updating message")
	err := m.repository.Update(dto)
	if err != nil {
		log.Error("error with updating message", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message updated")
	return nil
}

func (m *Service) Delete(id uuid.UUID) error {
	const op = "services.messenger.Delete"
	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("deleting message")
	err := m.repository.Delete(id)
	if err != nil {
		log.Error("error with deleting message", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("message deleted")
	return nil
}
