package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"messenger/internal/domain"
	"net/http"
)

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("upgrade websocket error", "error", err)
		return
	}
	h.log.Info(fmt.Sprintf("Client with address %s connected", conn.RemoteAddr().String()))
	go h.Conn(conn)
}

func (h *Handler) Conn(conn *websocket.Conn) {
	for {
		msg := new(domain.MessageAdd)
		if err := conn.ReadJSON(msg); err != nil {
			h.log.Warn("Error with reading from WebSocket: ", err)
			break
		}

		err := h.messengerService.Add(*msg)
		if err != nil {
			h.log.Warn("Error with adding message to Messenger: ", err)

		}
	}
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) writeToClientsBroadcast() {
	//for msg := range h.broadcast {
	//	h.mu.RLock()
	//	for client := range h.wsClients {
	//		go func(client *websocket.Conn) {
	//			if err := client.WriteJSON(msg); err != nil {
	//				h.log.Warn("Error with client write to clients: ", err)
	//			}
	//		}(client)
	//	}
	//	h.mu.RUnlock()
	//}
}
