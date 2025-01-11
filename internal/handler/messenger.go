package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"messenger/internal/domain"
	"net"
	"net/http"
	"time"
)

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("upgrade websocket error", "error", err)
		return
	}
	h.log.Info(fmt.Sprintf("Client with address %s connected", conn.RemoteAddr().String()))

	h.mu.Lock()
	h.wsClients[conn] = struct{}{}
	h.mu.Unlock()

	go h.readFromClient(conn)
}

func (h *Handler) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(domain.WSMessage)
		if err := conn.ReadJSON(msg); err != nil {
			h.log.Warn("Error with reading from WebSocket: ", err)
			break
		}
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			h.log.Warn("Error with address split: ", err)
		}

		msg.IPAddress = host
		msg.Time = time.Now().Format("15:04")
		h.broadcast <- msg
	}
	h.mu.Lock()
	delete(h.wsClients, conn)
	h.mu.Unlock()
}

func (h *Handler) writeToClientsBroadcast() {
	for msg := range h.broadcast {
		h.mu.RLock()
		for client := range h.wsClients {
			go func(client *websocket.Conn) {
				if err := client.WriteJSON(msg); err != nil {
					h.log.Warn("Error with client write to clients: ", err)
				}
			}(client)
		}
		h.mu.RUnlock()
	}
}
