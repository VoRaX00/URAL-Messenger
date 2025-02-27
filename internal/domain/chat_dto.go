package domain

import (
	"github.com/google/uuid"
	"messenger/internal/domain/models"
)

type AddChat struct {
	PersonIds []uuid.UUID `json:"personIds"`
	Name      string      `json:"name"`
}

type UpdateChat struct {
	Name string `json:"name"`
}

type GetChat struct {
	Name        string         `json:"name"`
	LastMessage models.Message `json:"lastMessage"`
}
