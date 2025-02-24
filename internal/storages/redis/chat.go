package redis

import (
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"messenger/internal/domain/models"
)

type ChatRepository struct {
	db *redis.Client
}

func NewChatRepository(client *redis.Client) *ChatRepository {
	return &ChatRepository{
		db: client,
	}
}

func (c *ChatRepository) Add(chat models.Chat, personIds []uuid.UUID) error {
	panic("implement me")
}
