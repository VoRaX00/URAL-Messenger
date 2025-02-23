package domain

import "github.com/google/uuid"

type AddChat struct {
	PersonIds []uuid.UUID `json:"personIds"`
	Name      string      `json:"name"`
}

type UpdateChat struct {
	Name string `json:"name"`
}
