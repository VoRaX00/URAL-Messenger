package chat

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
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
		addChat domain.AddChat
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
		input           args
		args            args
		expectedId      uuid.UUID
		expectedError   error
		mockParam       models.Chat
		mockReturnId    uuid.UUID
		mockReturnError error
	}{
		{
			name: "Успешное добавление",
			input: args{
				addChat: domain.AddChat{
					PersonIds: personsIds,
					Name:      "TestChat",
				},
			},
			expectedId:      chatId,
			expectedError:   nil,
			mockReturnId:    chatId,
			mockReturnError: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepository.On("Add", mock.AnythingOfType("models.Chat"),
				mock.AnythingOfType("[]uuid.UUID")).Return(c.mockReturnId, c.mockReturnError).Once()

			id, err := service.Add(c.input.addChat)
			require.Equal(t, c.expectedId, id)
			require.Equal(t, c.expectedError, err)
		})
	}
}

func TestService_AddNewUser(t *testing.T) {
	type args struct {
		chatId, userId uuid.UUID
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockRepository := mocks.NewRepository(t)
	mockCacheRepository := mocks.NewCacheRepository(t)

	service := NewChatService(slog.New(logHandler), mockRepository, mockCacheRepository)

	cases := []struct {
		name            string
		input           args
		mockReturnError error
		expectedError   error
	}{
		{
			name: "Успешное добавление пользователя в чат",
			input: args{
				chatId: uuid.New(),
				userId: uuid.New(),
			},
			mockReturnError: nil,
			expectedError:   nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("AddNewUser", mock.AnythingOfType("uuid.UUID"),
				mock.AnythingOfType("uuid.UUID")).Return(tt.mockReturnError).Once()

			err := service.AddNewUser(tt.input.chatId, tt.input.userId)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestService_RemoveUser(t *testing.T) {
	type args struct {
		chatId, userId uuid.UUID
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockRepository := mocks.NewRepository(t)
	mockCacheRepository := mocks.NewCacheRepository(t)
	service := NewChatService(slog.New(logHandler), mockRepository, mockCacheRepository)

	cases := []struct {
		name            string
		input           args
		mockReturnError error
		expectedError   error
	}{
		{
			name: "Успешное удаление пользователя из чата",
			input: args{
				chatId: uuid.New(),
				userId: uuid.New(),
			},
			mockReturnError: nil,
			expectedError:   nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("RemoveUser", mock.AnythingOfType("uuid.UUID"),
				mock.AnythingOfType("uuid.UUID")).Return(tt.mockReturnError).Once()

			err := service.RemoveUser(tt.input.chatId, tt.input.userId)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestService_GetUserChats(t *testing.T) {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockRepository := mocks.NewRepository(t)
	mockCacheRepository := mocks.NewCacheRepository(t)

	service := NewChatService(slog.New(logHandler), mockRepository, mockCacheRepository)

	chatId := uuid.New()

	cases := []struct {
		name               string
		mockReturnChatsIds []uuid.UUID
		mockReturnError    error
		expectedIds        []uuid.UUID
		expectedError      error
	}{
		{
			name:               "Успешное получение чата",
			mockReturnChatsIds: []uuid.UUID{chatId},
			mockReturnError:    nil,
			expectedIds:        []uuid.UUID{chatId},
			expectedError:      nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("GetUserChats", mock.AnythingOfType("uuid.UUID")).
				Return(tt.mockReturnChatsIds, tt.mockReturnError).Once()

			ids, err := service.GetUserChats(chatId)
			require.Equal(t, tt.expectedIds, ids)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestService_Update(t *testing.T) {
	type args struct {
		chat models.Chat
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	mockRepository := mocks.NewRepository(t)
	mockCacheRepository := mocks.NewCacheRepository(t)

	service := NewChatService(slog.New(logHandler), mockRepository, mockCacheRepository)

	cases := []struct {
		name            string
		input           args
		mockReturnError error
		expectedError   error
	}{
		{
			name: "Успешное обновлениее чата",
			input: args{
				chat: models.Chat{
					Id:   uuid.New(),
					Name: "TestChat",
				},
			},
			mockReturnError: nil,
			expectedError:   nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("Update", tt.input.chat).Return(tt.mockReturnError).Once()
			err := service.Update(tt.input.chat)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestService_Delete(t *testing.T) {
	type args struct {
		chatId, userId uuid.UUID
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	mockRepository := mocks.NewRepository(t)
	mockCacheRepository := mocks.NewCacheRepository(t)

	service := NewChatService(slog.New(logHandler), mockRepository, mockCacheRepository)
	cases := []struct {
		name            string
		input           args
		mockReturnError error
		expectedError   error
	}{
		{
			name: "Успешное удаление чата",
			input: args{
				chatId: uuid.New(),
				userId: uuid.New(),
			},
			mockReturnError: nil,
			expectedError:   nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("Delete", tt.input.chatId, tt.input.userId).Return(tt.mockReturnError).Once()
			err := service.Delete(tt.input.chatId, tt.input.userId)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestService_GetInfoUserChats(t *testing.T) {
	type args struct {
		userId      uuid.UUID
		page, count uint
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	mockRepository := mocks.NewRepository(t)
	mockCacheRepository := mocks.NewCacheRepository(t)
	service := NewChatService(slog.New(logHandler), mockRepository, mockCacheRepository)

	cases := []struct {
		name                    string
		input                   args
		mockReturnChatsId       []uuid.UUID
		mockReturnChatsIdError  error
		mockReturnInfoChats     []domain.GetChat
		mockReturnInfoChatError []error
		expectedChats           []domain.GetChat
		expectedError           error
	}{
		{
			name: "Успешное получение информации о чате на первой странице",
			input: args{
				userId: uuid.New(),
				page:   1,
				count:  1,
			},
			mockReturnChatsId:      []uuid.UUID{uuid.New()},
			mockReturnChatsIdError: nil,
			mockReturnInfoChats: []domain.GetChat{
				{
					Name: "Chat1",
					LastMessage: models.Message{
						MessageText: "Message1",
					},
				},
			},
			mockReturnInfoChatError: []error{
				nil,
			},
			expectedChats: []domain.GetChat{
				{
					Name: "Chat1",
					LastMessage: models.Message{
						MessageText: "Message1",
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.On("GetChatIds", tt.input.userId, tt.input.page,
				tt.input.count).Return(tt.mockReturnChatsId, tt.mockReturnChatsIdError).Once()

			for i, id := range tt.mockReturnChatsId {
				mockRepository.On("GetInfoChat", id).
					Return(tt.mockReturnInfoChats[i], tt.mockReturnInfoChatError[i]).Once()
			}

			chats, err := service.GetInfoUserChats(tt.input.userId, tt.input.page, tt.input.count)
			require.Equal(t, tt.expectedChats, chats)
			require.Equal(t, tt.expectedError, err)
		})
	}
}
