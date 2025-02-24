package chat

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
)

//go:generate mockery --name=CacheRepository --output=./mocks --case=underscore
type CacheRepository interface {
	Add(chat models.Chat, personIds []uuid.UUID) error
}

//go:generate mockery --name=Repository --output=./mocks --case=underscore
type Repository interface {
	Add(chat models.Chat, personIds []uuid.UUID) (uuid.UUID, error)
	AddNewUser(chatId uuid.UUID, userId uuid.UUID) error
	RemoveUser(chatId uuid.UUID, userId uuid.UUID) error
	GetUserChats(userId uuid.UUID) ([]uuid.UUID, error)
	Update(chat models.Chat) error
	Delete(chatId uuid.UUID) error
}

type Service struct {
	log             *slog.Logger
	repository      Repository
	cacheRepository CacheRepository
}

func NewChat(log *slog.Logger, repository Repository, cacheRepository CacheRepository) *Service {
	return &Service{
		log:             log,
		repository:      repository,
		cacheRepository: cacheRepository,
	}
}

func (c *Service) Add(chat domain.AddChat) (uuid.UUID, error) {
	panic("implement me")
}

func (c *Service) AddNewUser(chatId uuid.UUID, userId uuid.UUID) error {
	panic("implement me")
}

func (c *Service) RemoveUser(chatId uuid.UUID, userId uuid.UUID) error {
	panic("implement me")
}

func (c *Service) GetUserChats(userId uuid.UUID) ([]uuid.UUID, error) {
	const op = "services.messenger.GetUserChats"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("getting chats for user")
	chats, err := c.repository.GetUserChats(userId)
	if err != nil {
		log.Error("error with getting chats for user:", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("chats received")
	return chats, nil
}

func (c *Service) Update(chat models.Chat) error {
	panic("implement me")
}

func (c *Service) Delete(chatId uuid.UUID) error {
	panic("implement me")
}
