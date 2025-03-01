package postgres

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"messenger/internal/domain/models"
)

var ErrNotFound = errors.New("not found")

type ChatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (c *ChatRepository) Add(chat models.Chat, personIds []uuid.UUID) (uuid.UUID, error) {
	const op = "postgres.ChatRepository.Add"

	tx, err := c.db.Beginx()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	query := `INSERT INTO chats (id, name) VALUES ($1, $2)`
	_, err = tx.Exec(query, chat.Id, chat.Name)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	query = `INSERT INTO chats_persons (chat_id, person_id) VALUES ($1, $2)`
	for _, personId := range personIds {
		_, err = tx.Exec(query, chat.Id, personId)
		if err != nil {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return chat.Id, nil
}

func (c *ChatRepository) AddNewUser(chatId uuid.UUID, userId uuid.UUID) error {
	const op = "postgres.ChatRepository.AddNewUser"
	tx, err := c.db.Beginx()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	query := `SELECT EXISTS (SELECT 1 FROM chats WHERE id = $1)`

	var exists bool
	err = tx.QueryRow(query, chatId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = `INSERT INTO chats_persons (chat_id, person_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, chatId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (c *ChatRepository) RemoveUser(chatId uuid.UUID, userId uuid.UUID) error {
	const op = "postgres.ChatRepository.RemoveUser"
	tx, err := c.db.Beginx()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	query := `DELETE FROM chats_persons WHERE chat_id = $1 AND person_id = $2`
	_, err = tx.Exec(query, chatId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
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

func (c *ChatRepository) Update(chat models.Chat) error {
	const op = "postgres.ChatRepository.Update"
	tx, err := c.db.Beginx()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	query := `SELECT EXISTS (SELECT 1 FROM chats WHERE id = $1)`
	var exists bool
	err = tx.QueryRow(query, chat.Id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = `UPDATE chats SET name = $1 WHERE id = $2`
	_, err = tx.Exec(query, chat.Name, chat.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (c *ChatRepository) Delete(chatId, userId uuid.UUID) error {
	const op = "postgres.ChatRepository.Delete"
	tx, err := c.db.Beginx()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM public.chats_persons WHERE chat_id = $1 AND person_id = $2)`
	err = tx.QueryRow(query, chatId, userId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	query = `DELETE FROM chats WHERE id = $1`
	_, err = tx.Exec(query, chatId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = `DELETE FROM chats_persons WHERE chat_id = $1`
	_, err = tx.Exec(query, chatId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
