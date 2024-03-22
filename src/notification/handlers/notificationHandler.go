package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type NotificationHandler struct {
	upgrader websocket.Upgrader
}

func NewNotificationHandler() NotificationHandler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return NotificationHandler{
		upgrader: upgrader,
	}
}

func (nh NotificationHandler) WebSocketHandler(c *gin.Context) {
	ws, err := nh.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		ws.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
		time.Sleep(time.Minute)
	}
}
