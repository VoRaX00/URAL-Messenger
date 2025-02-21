package handler

import (
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
	"testing"
)

func TestWsConnectionHandler_Success(t *testing.T) {
	mockMessengerService := mocks.NewMessengerService(t)
	mockMessengerService.On("Add", mock.AnythingOfType("domain.MessageAdd")).Return(models.Message{}, nil)

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	h := NewHandler(slog.New(logHandler), mockMessengerService)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.wsHandler(w, r)
	}))

	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	msg := domain.MessageAdd{
		PersonId: uuid.New(),
		ChatId:   uuid.New(),
		Message:  "Hello tests",
	}
	err = conn.WriteJSON(msg)
	require.NoError(t, err)
}

func TestWsHandler_Fail(t *testing.T) {}
