package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/internal/services/mocks"
	"os"
	"testing"
	"time"
)

func TestMessenger_Add(t *testing.T) {
	type args struct {
		message models.Message
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockMessengerRepo := mocks.NewMessengerRepo(t)
	mockMessengerCacheRepo := mocks.NewMessengerCacheRepo(t)

	service := &Messenger{
		log:        slog.New(logHandler),
		cache:      mockMessengerCacheRepo,
		repository: mockMessengerRepo,
	}

	personId := uuid.New()
	chatId := uuid.New()
	msgId := uuid.New()
	textMsg := "TestAdd"
	timeNow := time.Now()

	cases := []struct {
		name              string
		Message           domain.MessageAdd
		mockArgument      args
		mockReturnMessage models.Message
		mockReturnError   error
		expectedMessage   models.Message
		expectedError     error
	}{
		{
			name: "Успешное добавление",
			Message: domain.MessageAdd{
				PersonId: personId,
				ChatId:   chatId,
				Message:  textMsg,
			},
			mockArgument: args{
				message: models.Message{
					MessageText: textMsg,
					Chat: models.Chat{
						Id: chatId,
					},
					PersonId: personId,
				},
			},
			mockReturnMessage: models.Message{
				Id:          msgId,
				MessageText: textMsg,
				Chat: models.Chat{
					Id: chatId,
				},
				PersonId:    personId,
				SendingTime: timeNow,
			},
			mockReturnError: nil,
			expectedMessage: models.Message{
				Id:          msgId,
				MessageText: textMsg,
				Chat: models.Chat{
					Id: chatId,
				},
				PersonId:    personId,
				SendingTime: timeNow,
			},
			expectedError: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			mockMessengerRepo.ExpectedCalls = nil
			mockMessengerCacheRepo.ExpectedCalls = nil

			mockMessengerRepo.On("Add", c.mockArgument.message).Return(c.mockReturnMessage, c.mockReturnError).Once()
			mockMessengerCacheRepo.On("Add", c.mockArgument.message).Return(c.mockReturnError).Once()

			msg, err := service.Add(c.Message)
			require.Equal(t, c.expectedMessage, msg)
			require.Equal(t, c.expectedError, err)
		})
	}
}
