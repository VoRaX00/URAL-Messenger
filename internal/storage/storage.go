package storage

import (
	"github.com/google/uuid"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
)

//type Storage interface {
//	MustConnect()
//	Connect() error
//	MustClose()
//	Close() error
//}

type MessengerCacheRepo interface {
	Add(message models.Message) error
	Delete(message models.Message) error
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
}

type MessengerRepo interface {
	Add(message models.Message) error
	GetByChat(chatId uuid.UUID) ([]models.Message, error)
	GetById(id uuid.UUID) (models.Message, error)
	Update(message domain.MessageUpdate) error
	Delete(id uuid.UUID) error
}
