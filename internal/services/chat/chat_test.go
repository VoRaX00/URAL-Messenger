package chat

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/internal/services/chat/mocks"
	"os"
	"testing"
)

func TestService_Add(t *testing.T) {
	type args struct {
		chat       models.Chat
		personsIds []uuid.UUID
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockCacheRepository := mocks.NewCacheRepository(t)
	mockRepository := mocks.NewRepository(t)

	service := NewChatService(slog.New(logHandler), mockRepository, mockCacheRepository)

	chatId := uuid.New()
	personsIds := []uuid.UUID{uuid.New(), uuid.New()}

	cases := []struct {
		name            string
		input           domain.AddChat
		args            args
		expectedId      uuid.UUID
		expectedError   error
		mockReturnId    uuid.UUID
		mockReturnError error
	}{
		{
			name: "Успешное добавление",
			input: domain.AddChat{
				PersonIds: personsIds,
				Name:      "TestChat",
			},
			args: args{
				chat: models.Chat{
					Id:   chatId,
					Name: "TestChat",
				},
				personsIds: personsIds,
			},
			expectedId:      chatId,
			expectedError:   nil,
			mockReturnId:    chatId,
			mockReturnError: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepository.On("Add", c.args.chat, c.args.personsIds).Return(c.mockReturnId, c.mockReturnError).Once()

			id, err := service.Add(c.input)
			require.Equal(t, c.expectedId, id)
			require.Equal(t, c.expectedError, err)
		})
	}
}
