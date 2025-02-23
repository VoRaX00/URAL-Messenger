package redis

import (
	"github.com/go-redis/redis"
	"messenger/internal/domain"
)

type ChatRepository struct {
	db *redis.Client
}

func NewChatRepository(client *redis.Client) *ChatRepository {
	return &ChatRepository{
		db: client,
	}
}

func (c *ChatRepository) Add(chat domain.AddChat) error {
	panic("implement me")
}
