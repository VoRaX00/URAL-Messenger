package redis

import (
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"messenger/internal/domain/models"
)

type CacheMessenger struct {
	db *redis.Client
}

func NewMessengerCacheRepo(db *redis.Client) *CacheMessenger {
	return &CacheMessenger{
		db: db,
	}
}

func (m *CacheMessenger) Add(message models.Message) error {
	panic("implement me")
}

func (m *CacheMessenger) Delete(message models.Message) error {
	panic("implement me")
}

func (m *CacheMessenger) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	panic("implement me")
}
