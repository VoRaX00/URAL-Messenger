package redis

import (
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"messenger/internal/domain/models"
)

type MessageRepository struct {
	db *redis.Client
}

func NewMessageRepository(db *redis.Client) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (m *MessageRepository) Add(message models.Message) error {
	panic("implement me")
}

func (m *MessageRepository) Delete(message models.Message) error {
	panic("implement me")
}

func (m *MessageRepository) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	panic("implement me")
}
