package chat

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

func NewChatService(log *slog.Logger, repository Repository, cacheRepository CacheRepository) *Service {
	return &Service{
		log:             log,
		repository:      repository,
		cacheRepository: cacheRepository,
	}
}

func (c *Service) Add(addChat domain.AddChat) (uuid.UUID, error) {
	const op = "service.chat.Add"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("mapping addChat to Chat")
	chat := mapper.AddChatToChat(addChat)
	chat.Id = uuid.New()
	log.Info("successfully mapped addChat to Chat")

	log.Info("adding new chat")
	id, err := c.repository.Add(chat, addChat.PersonIds)
	if err != nil {
		log.Error("Error with adding chat to repository:", slog.String("err", err.Error()))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully added new chat to repository")

	return id, nil
}

func (c *Service) AddNewUser(chatId uuid.UUID, userId uuid.UUID) error {
	const op = "service.chat.AddNewUser"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("adding new user to chat")
	err := c.repository.AddNewUser(chatId, userId)
	if err != nil {
		log.Error("Error with adding new user to repository:", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully added new user to repository")

	return nil
}

func (c *Service) RemoveUser(chatId uuid.UUID, userId uuid.UUID) error {
	const op = "service.chat.RemoveUser"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("removing user from chat")
	err := c.repository.RemoveUser(chatId, userId)
	if err != nil {
		log.Error("Error with removing user from repository:", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully removed user from repository")

	return nil
}

func (c *Service) GetInfoUserChats(userId uuid.UUID, page, count uint) ([]domain.GetChat, error) {
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
	const op = "services.messenger.Update"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("updating chat")
	err := c.repository.Update(chat)
	if err != nil {
		log.Error("error with updating chat:", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully updated chat")

	return nil
}

func (c *Service) Delete(chatId, userId uuid.UUID) error {
	const op = "services.messenger.Delete"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("deleting chat")
	err := c.repository.Delete(chatId)
	if err != nil {
		log.Error("error with deleting chat:", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully deleted chat")

	return nil
}
