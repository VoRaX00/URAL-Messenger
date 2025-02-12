package models

import "github.com/google/uuid"

type File struct {
	Id        uuid.UUID `json:"id" db:"id"`
	File      byte      `json:"file" db:"file"`
	MessageId uuid.UUID `json:"message_id" db:"message_id"`
}
