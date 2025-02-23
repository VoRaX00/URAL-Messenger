package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"messenger/internal/domain/models"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (m *MessageRepository) Add(message models.Message) (models.Message, error) {
	const op = "MessengerRepo.Add"

	tx, err := m.db.Beginx()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err != nil {
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	chatExists, err := m.checkExistsChat(tx, message.Chat.Id)
	if err != nil {
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	if !chatExists {
		err = m.createChat(tx, message.Chat)
		if err != nil {
			return models.Message{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	query := `INSERT INTO messages (id, message, person_id, chat_id, sending_time) VALUES ($1, $2, $3, $4, $5)`
	_, err = m.db.Exec(query, message)
	if err != nil {
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	msg, err := m.GetById(message.Id)
	if err != nil {
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	return msg, nil
}

func (m *MessageRepository) checkExistsChat(tx *sqlx.Tx, chatID uuid.UUID) (bool, error) {
	const op = "MessengerRepo.checkExistsChat"
	query := `SELECT EXISTS (SELECT 1 FROM chats WHERE id = $1)`

	var exists bool
	if err := tx.QueryRow(query, chatID).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		return false, nil
	}
	return true, nil
}

func (m *MessageRepository) createChat(tx *sqlx.Tx, chat models.Chat, personsId ...uuid.UUID) error {
	const op = "MessengerRepo.CreateChat"

	query := `INSERT INTO chats (id, name) VALUES ($1, $2)`
	_, err := tx.Exec(query, chat.Id, chat.Name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = `INSERT INTO chats_persons (chat_id, person_id) VALUES ($1, $2)`
	for _, p := range personsId {
		_, err = tx.Exec(query, p, chat.Id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (m *MessageRepository) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	const op = `MessengerRepo.GetByChat`
	query := `SELECT id, message, person_id, chat_id, sending_time FROM messages WHERE chat_id = $1 AND status <> $2`

	var messages []models.Message
	err := m.db.Select(messages, query, chatId, "deleted")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return messages, nil
}

func (m *MessageRepository) GetById(id uuid.UUID) (models.Message, error) {
	const op = `MessengerRepo.GetById`
	query := `SELECT id, message, person_id, chat_id, sending_time FROM messages WHERE id = $1`

	var message models.Message
	err := m.db.Get(&message, query, id)
	if err != nil {
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}
	return message, nil
}

func (m *MessageRepository) GetUserChats(userId uuid.UUID) ([]uuid.UUID, error) {
	const op = `MessengerRepo.GetUserChats`
	query := `SELECT chat_id FROM chats_persons WHERE person_id = $1`
	var chats []uuid.UUID
	err := m.db.Select(&chats, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return chats, nil
}

func (m *MessageRepository) Update(message models.Message) error {
	const op = `MessengerRepo.Update`
	query := `UPDATE messages SET message=$1, status=$2 WHERE id = $3`
	_, err := m.db.Exec(query, message.MessageText, message, message.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (m *MessageRepository) Delete(id uuid.UUID) error {
	const op = `MessengerRepo.Delete`
	query := `UPDATE messages SET status = $1 WHERE id = $2`
	_, err := m.db.Exec(query, "deleted", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
