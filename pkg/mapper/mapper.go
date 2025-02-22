package mapper

import (
	"messenger/internal/domain"
	"messenger/internal/domain/models"
)

func MessageAddToMessage(message domain.MessageAdd) models.Message {
	return models.Message{
		PersonId: message.PersonId,
		Chat: models.Chat{
			Id: message.ChatId,
		},
		MessageText: message.Message,
	}
}

func MessageUpdateToMessage(message domain.MessageUpdate) models.Message {
	return models.Message{
		Id:          message.Id,
		MessageText: message.Message,
		Status:      message.Status,
	}
}
