package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"messenger/internal/domain"
	"messenger/internal/domain/models"
	"messenger/internal/handler/mocks"
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

func TestWsHandler_Fail(t *testing.T) {}
