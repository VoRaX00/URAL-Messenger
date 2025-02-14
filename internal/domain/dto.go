package domain

import "github.com/google/uuid"

type MessageAdd struct {
	PersonId uuid.UUID `json:"personId"`
	ChatId   uuid.UUID `json:"chatId"`
	Message  string    `json:"message"`
}

type MessageUpdate struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
	Status  string    `json:"status"`
}
