package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id          uuid.UUID `json:"id"`
	PersonId    uuid.UUID `json:"person_id" db:"person_id"`
	Chat        Chat      `json:"chat" db:"chat"`
	MessageText string    `json:"message" db:"message"`
	SendingTime time.Time `json:"time" db:"time"`
}
