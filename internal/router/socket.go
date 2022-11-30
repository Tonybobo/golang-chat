package router

import (
	"github.com/tonybobo/go-chat/internal/server"
	"github.com/tonybobo/go-chat/pkg/global/log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(c *gin.Context) {
	user := c.Query("username")
	if user == "" {
		return
	}

	log.Logger.Info("New User to the Server", log.String("User: ", user))

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logger.Error("Error at Websocket", log.Any("Error", err))
		return
	}

	client := &server.Client{
		Conn: ws,
		Send: make(chan []byte),
		User: user,
	}

	server.WebSocketServer.Online <- client

	go client.Receiver()
	go client.Writer()
}
