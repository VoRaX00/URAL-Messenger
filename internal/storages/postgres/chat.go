package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"messenger/internal/domain"
)

type ChatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (c *ChatRepository) Add(chat domain.AddChat) (uuid.UUID, error) {
	panic("implement me")
}

func (c *ChatRepository) AddNewUser(chatId uuid.UUID, userId uuid.UUID) error {
	panic("implement me")
}

func (c *ChatRepository) RemoveUser(chatId uuid.UUID, userId uuid.UUID) error {
	panic("implement me")
}

func (c *ChatRepository) GetUserChats(userId uuid.UUID) ([]uuid.UUID, error) {
	const op = `ChatRepository.GetUserChats`
	query := `SELECT chat_id FROM chats_persons WHERE person_id = $1`

	var chats []uuid.UUID
	err := c.db.Select(&chats, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return chats, nil
}

func (c *ChatRepository) Update(chatId uuid.UUID) error {
	panic("implement me")
}

func (c *ChatRepository) Delete(chatId uuid.UUID) error {
	panic("implement me")
}
