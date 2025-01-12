package handlers

import (
	ws "bonchDvach/pkg/websockets"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var wsHub = ws.NewHub()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsHub.RegisterClient(conn)
}

func InitWebSocketHub() {
	go wsHub.Run()
}
