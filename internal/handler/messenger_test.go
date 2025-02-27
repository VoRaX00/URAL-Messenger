package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/internal/handler/mocks"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWsConnection(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	userId := uuid.New()

	mockMessengerService := mocks.NewMessageService(t)
	mockChatService := mocks.NewChatService(t)
	mockChatService.On("GetUserChats", mock.AnythingOfType("uuid.UUID")).Return([]uuid.UUID{}, nil).
		Run(func(args mock.Arguments) {
			wg.Done()
		})

	h := NewHandler(slog.New(logHandler), mockMessengerService, mockChatService)
	h.InitRoutes()

	server := httptest.NewServer(h)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + fmt.Sprintf("/ws?user_id=%v", userId)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	conn.Close()

	wg.Wait()
	mockMessengerService.AssertExpectations(t)
}

func TestWsWrite_Success(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	person1 := uuid.New()
	textMsg := "Hello tests"
	chatId := uuid.New()

	mockMessengerService := mocks.NewMessageService(t)
	mockMessengerService.On("Add", mock.AnythingOfType("domain.MessageAdd")).Return(models.Message{
		PersonId:    person1,
		MessageText: textMsg,
		Chat: models.Chat{
			Id: chatId,
		},
	}, nil).Run(func(args mock.Arguments) {
		wg.Done()
	})

	mockChatService := mocks.NewChatService(t)
	mockChatService.On("GetUserChats", mock.AnythingOfType("uuid.UUID")).Return([]uuid.UUID{
		chatId,
	}, nil)

	h := NewHandler(slog.New(logHandler), mockMessengerService, mockChatService)
	h.InitRoutes()

	server := httptest.NewServer(h)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + fmt.Sprintf("/ws?user_id=%v", person1)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	msg := domain.MessageAdd{
		PersonId: person1,
		ChatId:   chatId,
		Message:  textMsg,
	}
	err = conn.WriteJSON(msg)
	require.NoError(t, err)

	wg.Wait()
	mockMessengerService.AssertExpectations(t)
}

func TestWsWriteReadForOneChat_Success(t *testing.T) {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	person1 := uuid.New()
	person2 := uuid.New()
	textMsg := "Hello tests"
	chatId := uuid.New()

	mockMessengerService := mocks.NewMessageService(t)
	mockMessengerService.On("Add", mock.AnythingOfType("domain.MessageAdd")).Return(models.Message{
		PersonId:    person1,
		MessageText: textMsg,
		Chat: models.Chat{
			Id: chatId,
		},
	}, nil)

	mockChatService := mocks.NewChatService(t)
	mockChatService.On("GetUserChats", mock.AnythingOfType("uuid.UUID")).Return([]uuid.UUID{
		chatId,
	}, nil)

	h := NewHandler(slog.New(logHandler), mockMessengerService, mockChatService)
	h.InitRoutes()

	server := httptest.NewServer(h)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + fmt.Sprintf("/ws?user_id=%v", person1)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	wsURL2 := "ws" + strings.TrimPrefix(server.URL, "http") + fmt.Sprintf("/ws?user_id=%v", person2)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	require.NoError(t, err)
	defer conn2.Close()

	msg := domain.MessageAdd{
		PersonId: person1,
		ChatId:   chatId,
		Message:  textMsg,
	}
	err = conn.WriteJSON(msg)
	require.NoError(t, err)

	var readMsg models.Message
	_ = conn2.SetReadDeadline(time.Now().Add(time.Second * 5))
	err = conn2.ReadJSON(&readMsg)
	require.NoError(t, err)
	require.Equal(t, textMsg, readMsg.MessageText)
}

func TestWsGetInfoUserChats(t *testing.T) {
	type args struct {
		userId      uuid.UUID
		page, count uint
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	user1 := uuid.New()

	mockMessengerService := mocks.NewMessageService(t)
	mockChatService := mocks.NewChatService(t)

	h := NewHandler(slog.New(logHandler), mockMessengerService, mockChatService)
	h.InitRoutes()

	cases := []struct {
		name            string
		input           args
		mockChatsReturn []domain.GetChat
		mockChatsError  error
		expectedChats   []domain.GetChat
		expectedStatus  int
	}{
		{
			name: "Возвращение одного чата",
			input: args{
				userId: user1,
				page:   1,
				count:  1,
			},
			mockChatsReturn: []domain.GetChat{
				{
					Name: "34 сквад",
					LastMessage: models.Message{
						PersonId:    user1,
						MessageText: "Тест",
					},
				},
			},
			mockChatsError: nil,
			expectedChats: []domain.GetChat{
				{
					Name: "34 сквад",
					LastMessage: models.Message{
						PersonId:    user1,
						MessageText: "Тест",
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	server := httptest.NewServer(h)
	defer server.Close()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockChatService.On("GetInfoUserChats", tt.input.userId, tt.input.page, tt.input.count).Return(tt.mockChatsReturn, tt.mockChatsError)
			url := fmt.Sprintf("%s/chat/info?userId=%v&page=%d&count=%d", server.URL, tt.input.userId, tt.input.page, tt.input.count)
			resp, err := http.Get(url)
			require.NoError(t, err)

			defer resp.Body.Close()
			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			var chats []domain.GetChat
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&chats)
			if err != nil {
				t.Error(err)
			}

			require.Equal(t, tt.expectedChats, chats)
		})
	}
}

func TestWsHandler_Fail(t *testing.T) {}
