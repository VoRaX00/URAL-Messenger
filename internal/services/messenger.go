package services

import (
	"github.com/google/uuid"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/internal/storage"
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
	panic("implement me")
}

func (m *Messenger) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	panic("implement me")
}

func (m *Messenger) GetById(id uuid.UUID) (models.Message, error) {
	panic("implement me")
}

func (m *Messenger) Update(id uuid.UUID, message string) error {
	panic("implement me")
}

func (m *Messenger) Delete(id uuid.UUID) error {
	panic("implement me")
}
