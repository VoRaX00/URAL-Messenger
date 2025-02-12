package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"messenger/internal/domain/models"
)

type Messenger struct {
	db *sqlx.DB
}

func NewMessengerRepo(db *sqlx.DB) *Messenger {
	return &Messenger{
		db: db,
	}
}

func (p *Messenger) Add(message models.Message) error {
	panic("implement me")
}

func (p *Messenger) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	panic("implement me")
}

func (p *Messenger) GetById(id uuid.UUID) (models.Message, error) {
	panic("implement me")
}

func (p *Messenger) Update(id uuid.UUID, message string) error {
	panic("implement me")
}

func (p *Messenger) Delete(id uuid.UUID) error {
	panic("implement me")
}
