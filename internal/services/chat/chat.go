package chat

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/pkg/mapper"
	"sync"
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
	GetChatIds(userId uuid.UUID, offset, limit uint) ([]uuid.UUID, error)
	GetInfoChat(chatId uuid.UUID) (domain.GetChat, error)
	Update(chat models.Chat) error
	Delete(chatId, userId uuid.UUID) error
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

// TODO: иногда падает deadlock, разобраться и понять почему так происходит
func (c *Service) GetInfoUserChats(userId uuid.UUID, page, count uint) ([]domain.GetChat, error) {
	const op = "services.messenger.GetInfoUserChats"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("getting user's chats")
	chatsIds, err := c.repository.GetChatIds(userId, page, count)
	if err != nil {
		log.Error("Error with getting user's chats:", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chats := make([]domain.GetChat, len(chatsIds))

	wg := sync.WaitGroup{}
	quit := make(chan error)
	mu := sync.Mutex{}

	for _, chatId := range chatsIds {
		wg.Add(1)
		go func() {
			defer wg.Done()

			info, err := c.repository.GetInfoChat(chatId)

			mu.Lock()
			chats = append(chats, info)
			if err != nil {
				quit <- err
			}
			mu.Unlock()
		}()
	}

	select {
	case err = <-quit:
		log.Error("Error with getting user's chats:", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	default:
		wg.Wait()
	}

	log.Info("successfully got user's chats")

	return chats, nil
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
	err := c.repository.Delete(chatId, userId)
	if err != nil {
		log.Error("error with deleting chat:", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully deleted chat")

	return nil
}
