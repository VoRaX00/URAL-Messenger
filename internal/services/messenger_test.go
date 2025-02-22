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

func TestMessenger_GetByChat(t *testing.T) {
	type args struct {
		chatId uuid.UUID
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

	chatId := uuid.New()

	cases := []struct {
		name               string
		chatId             uuid.UUID
		mockArgument       args
		mockReturnMessages []models.Message
		mockReturnError    error
		expectedMessages   []models.Message
		expectedError      error
	}{
		{
			name:   "Успешное получение чатов",
			chatId: chatId,
			mockArgument: args{
				chatId: chatId,
			},
			mockReturnMessages: []models.Message{},
			mockReturnError:    nil,
			expectedMessages:   []models.Message{},
			expectedError:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			mockMessengerRepo.ExpectedCalls = nil
			mockMessengerCacheRepo.ExpectedCalls = nil

			mockMessengerRepo.On("GetByChat", c.mockArgument.chatId).Return(c.mockReturnMessages, c.mockReturnError).Once()
			messages, err := service.GetByChat(c.chatId)
			require.Equal(t, c.expectedMessages, messages)
			require.Equal(t, c.expectedError, err)
		})
	}
}

func TestMessenger_GetById(t *testing.T) {
	type args struct {
		msgId uuid.UUID
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockMessengerRepo := mocks.NewMessengerRepo(t)

	msgId := uuid.New()

	service := &Messenger{
		log:        slog.New(logHandler),
		cache:      nil,
		repository: mockMessengerRepo,
	}

	cases := []struct {
		name              string
		msgId             uuid.UUID
		args              args
		mockReturnMessage models.Message
		mockReturnError   error
		expectedMessage   models.Message
		expectedError     error
	}{
		{
			name:  "Успешное получение сообщения",
			msgId: msgId,
			args: args{
				msgId: msgId,
			},
			mockReturnMessage: models.Message{
				Id: msgId,
			},
			mockReturnError: nil,
			expectedMessage: models.Message{
				Id: msgId,
			},
			expectedError: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockMessengerRepo.ExpectedCalls = nil
			mockMessengerRepo.On("GetById", c.args.msgId).Return(c.mockReturnMessage, c.mockReturnError).Once()

			message, err := service.GetById(c.msgId)
			require.Equal(t, c.expectedMessage, message)
			require.Equal(t, c.expectedError, err)
		})
	}
}

func TestMessenger_Update(t *testing.T) {
	type args struct {
		msg models.Message
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockMessengerRepo := mocks.NewMessengerRepo(t)
	msgId := uuid.New()
	msg := domain.MessageUpdate{
		Id:      msgId,
		Message: "Updated msg",
		Status:  "not read",
	}

	service := &Messenger{
		log:        slog.New(logHandler),
		cache:      nil,
		repository: mockMessengerRepo,
	}

	cases := []struct {
		name            string
		msg             domain.MessageUpdate
		args            args
		mockReturnError error
		expectedError   error
	}{
		{
			name: "Успешное изменение сообщения",
			msg:  msg,
			args: args{
				msg: models.Message{
					Id:          msgId,
					MessageText: msg.Message,
					Status:      msg.Status,
				},
			},
			mockReturnError: nil,
			expectedError:   nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockMessengerRepo.ExpectedCalls = nil
			mockMessengerRepo.On("Update", c.args.msg).Return(c.mockReturnError).Once()

			err := service.Update(msg)
			require.Equal(t, c.expectedError, err)
		})
	}
}

func TestMessenger_Delete(t *testing.T) {
	type args struct {
		msgId uuid.UUID
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockMessengerRepo := mocks.NewMessengerRepo(t)
	msgId := uuid.New()

	service := &Messenger{
		log:        slog.New(logHandler),
		cache:      nil,
		repository: mockMessengerRepo,
	}

	cases := []struct {
		name            string
		input           uuid.UUID
		args            args
		mockReturnError error
		expectedError   error
	}{
		{
			name:  "Успешное удаление",
			input: msgId,
			args: args{
				msgId: msgId,
			},
			mockReturnError: nil,
			expectedError:   nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockMessengerRepo.ExpectedCalls = nil
			mockMessengerRepo.On("Delete", c.args.msgId).Return(c.mockReturnError).Once()
			err := service.Delete(c.input)
			require.Equal(t, c.expectedError, err)
		})
	}
}
